package bee

import (
	"context"
	"gitee.com/stuinfer/bee-api/common"
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/model"
	"gitee.com/stuinfer/bee-api/printer"
	"gitee.com/stuinfer/bee-api/service"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee"
	beeReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/model/bee/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/beeshop/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"text/template"
)

type BeePrinterService struct{}

// CreateBeePrinter 创建beePrinter表记录
// Author [piexlmax](https://github.com/piexlmax)
func (beePrinterService *BeePrinterService) CreateBeePrinter(beePrinter *bee.BeePrinter) (err error) {
	beePrinter.DateAdd = utils.NowPtr()
	beePrinter.DateUpdate = utils.NowPtr()
	if err = beePrinterService.checkTemplate(beePrinter.Template); err != nil {
		return err
	}
	return GetBeeDB().Transaction(func(tx *gorm.DB) error {
		if err := GetBeeDB().Create(beePrinter).Error; err != nil {
			return err
		}
		cfg := beePrinterService.printer2apiPrinter(beePrinter)
		p := printer.GetPrinter(cfg)
		if p == nil {
			err = errors.New("该品牌暂不支持")
			return err
		}
		return p.AddPrinter(cfg)
	})
}

func (beePrinterService *BeePrinterService) checkTemplate(tpl string) error {
	t := template.New("tmp")
	_, err := t.Parse(tpl)
	if err != nil {
		return err
	}
	return nil
}

func (beePrinterService *BeePrinterService) printer2apiPrinter(beePrinter *bee.BeePrinter) *model.BeePrinter {
	return &model.BeePrinter{
		BaseModel: common.BaseModel{
			Id:        cast.ToInt64(beePrinter.Id),
			UserId:    cast.ToInt64(beePrinter.UserId),
			IsDeleted: cast.ToBool(beePrinter.IsDeleted),
		},
		Appid:     beePrinter.Appid,
		AppSecret: beePrinter.AppSecret,
		Name:      beePrinter.Name,
		Key:       beePrinter.Key,
		Brand:     enum.PrinterBrand(cast.ToInt32(beePrinter.Brand)),
		Code:      beePrinter.Code,
		Condition: enum.PrinterCondition(cast.ToInt32(beePrinter.Condition)),
		Num:       cast.ToInt(beePrinter.Num),
		ShopId:    cast.ToInt64(beePrinter.ShopId),
		Template:  beePrinter.Template,
	}
}

// DeleteBeePrinter 删除beePrinter表记录
// Author [piexlmax](https://github.com/piexlmax)
func (beePrinterService *BeePrinterService) DeleteBeePrinter(id string, shopUserId int) (err error) {
	printerInfo, err := beePrinterService.GetBeePrinter(id, shopUserId)
	if err != nil {
		return err
	}
	cfg := beePrinterService.printer2apiPrinter(&printerInfo)
	err = printer.GetPrinter(cfg).DelPrinter(cfg, []string{printerInfo.Code})
	if err != nil {
		return err
	}
	err = GetBeeDB().Model(&bee.BeePrinter{}).Where("id = ?", id).Where("user_id = ?", shopUserId).
		Updates(map[string]interface{}{
			"is_deleted":  1,
			"date_delete": utils.NowPtr(),
		}).Error
	return err
}

// DeleteBeePrinterByIds 批量删除beePrinter表记录
// Author [piexlmax](https://github.com/piexlmax)
func (beePrinterService *BeePrinterService) DeleteBeePrinterByIds(ids []string, shopUserId int) (err error) {
	for _, id := range ids {
		if err = beePrinterService.DeleteBeePrinter(id, shopUserId); err != nil {
			return err
		}
	}
	return err
}

// UpdateBeePrinter 更新beePrinter表记录
// Author [piexlmax](https://github.com/piexlmax)
func (beePrinterService *BeePrinterService) UpdateBeePrinter(beePrinter bee.BeePrinter, shopUserId int) (err error) {
	beePrinter.DateUpdate = utils.NowPtr()
	if err = beePrinterService.checkTemplate(beePrinter.Template); err != nil {
		return err
	}

	var needReBind = false
	oldPrinterInfo, err := beePrinterService.GetBeePrinter(cast.ToString(beePrinter.Id), shopUserId)
	if err != nil {
		return err
	}
	if oldPrinterInfo.Code != beePrinter.Code || cast.ToInt(beePrinter.Brand) != cast.ToInt(oldPrinterInfo.Brand) ||
		beePrinter.Appid != oldPrinterInfo.Appid || beePrinter.AppSecret != oldPrinterInfo.AppSecret ||
		beePrinter.Key != oldPrinterInfo.Key || beePrinter.Name != oldPrinterInfo.Name {
		needReBind = true
	}
	// 先解绑再绑定
	if needReBind {
		cfg := beePrinterService.printer2apiPrinter(&oldPrinterInfo)
		p := printer.GetPrinter(cfg)
		if p != nil {
			if err := p.DelPrinter(cfg, []string{cfg.Code}); err != nil {
				return errors.Wrap(err, "移除旧打印机失败")
			} else {
				global.GVA_LOG.Info("移除旧的打印机成功")
			}
		} else {
			global.GVA_LOG.Info("旧的打印机不支持移除，跳过")
		}
	}

	return GetBeeDB().Transaction(func(tx *gorm.DB) error {
		if err := GetBeeDB().Model(&bee.BeePrinter{}).Where("id = ? and user_id = ?", beePrinter.Id, shopUserId).Updates(&beePrinter).Error; err != nil {
			return err
		}
		cfg := beePrinterService.printer2apiPrinter(&beePrinter)
		p := printer.GetPrinter(cfg)
		if p == nil {
			err = errors.New("该品牌暂不支持")
			return err
		}
		if needReBind {
			if err := p.AddPrinter(cfg); err != nil {
				return err
			}
		}
		return nil
	})
}

// TestBeePrinter 测试打印机配置是否正确
// Author [piexlmax](https://github.com/piexlmax)
func (beePrinterService *BeePrinterService) TestBeePrinter(beePrinter bee.BeePrinter, shopUserId int) (err error) {
	beePrinter.DateUpdate = utils.NowPtr()
	if err = beePrinterService.checkTemplate(beePrinter.Template); err != nil {
		return err
	}
	cfg := beePrinterService.printer2apiPrinter(&beePrinter)
	p := printer.GetPrinter(cfg)
	if p == nil {
		err = errors.New("该品牌暂不支持")
		return err
	}
	var needReBind = false
	if cfg.Id == 0 {
		needReBind = true
	} else {
		oldPrinterInfo, err := beePrinterService.GetBeePrinter(cast.ToString(beePrinter.Id), shopUserId)
		if err != nil {
			return err
		}
		if oldPrinterInfo.Code != beePrinter.Code || cast.ToInt(beePrinter.Brand) != cast.ToInt(oldPrinterInfo.Brand) ||
			beePrinter.Appid != oldPrinterInfo.Appid || beePrinter.AppSecret != oldPrinterInfo.AppSecret ||
			beePrinter.Key != oldPrinterInfo.Key || beePrinter.Name != oldPrinterInfo.Name {
			needReBind = true
		}
	}
	if needReBind {
		if err := p.AddPrinter(cfg); err != nil {
			if !errors.Is(err, printer.ErrAlreadyAdd) {
				return err
			}
		}
		defer func() {
			if err := p.DelPrinter(cfg, []string{cfg.Code}); err != nil {
				global.GVA_LOG.Error("移除打印机失败", zap.Error(err))
			}
		}()
	}
	return service.GetPrinterSrv().TestPrinter(context.Background(), cfg)
}

// GetBeePrinter 根据id获取beePrinter表记录
// Author [piexlmax](https://github.com/piexlmax)
func (beePrinterService *BeePrinterService) GetBeePrinter(id string, shopUserId int) (beePrinter bee.BeePrinter, err error) {
	err = GetBeeDB().Where("id = ? and user_id = ?", id, shopUserId).First(&beePrinter).Error
	return
}

// GetBeePrinterInfoList 分页获取beePrinter表记录
// Author [piexlmax](https://github.com/piexlmax)
func (beePrinterService *BeePrinterService) GetBeePrinterInfoList(info beeReq.BeePrinterSearch, shopUserId int) (list []bee.BeePrinter, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := GetBeeDB().Model(&bee.BeePrinter{})
	db = db.Where("user_id = ? and is_deleted = 0", shopUserId)
	var beePrinters []bee.BeePrinter
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&beePrinters).Error
	return beePrinters, total, err
}
