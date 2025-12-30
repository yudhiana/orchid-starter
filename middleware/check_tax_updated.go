package middleware

import (
	"orchid-starter/constants"
	"orchid-starter/internal/clients"
	"orchid-starter/internal/common"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/x/errors"
	"github.com/mataharibiz/sange/v2"
)

func CheckTaxUpdated(irisCtx iris.Context) *sange.Error {
	client := clients.NewClient()
	internalClient := client.InternalClient
	ctx := common.SetRequestContext(irisCtx.Request().Context(), irisCtx)

	var (
		companyID = common.GetCompanyIDFromContext(ctx)
		appOrigin = common.GetAppOriginFromContext(ctx)
	)

	if constants.IsSeller(appOrigin) {
		selected := `items { updateProductTaxStatus type }`
		result, errGet := internalClient.GetCompanyGQLDetail(ctx, companyID, selected)
		if errGet != nil {
			return sange.SetError(sange.Forbidden, errGet, "get gql company detail")
		}

		whiteListCompanyType := map[string]bool{
			strconv.FormatInt(constants.GetConstant("TYPE_VENDOR"), 10): true,
			strconv.FormatInt(constants.GetConstant("TYPE_BOTH"), 10):   true,
		}

		whiteListProductTaxStatus := map[string]iris.Map{
			strconv.FormatInt(constants.GetConstant("PRODUCT_TAX_NEED_UPDATE"), 10): {
				"error_code":    "PRODUCT_TAX_NEED_UPDATE",
				"error_message": "need update product tax",
			},
			strconv.FormatInt(constants.GetConstant("PRODUCT_TAX_INPROGRESS"), 10): {
				"error_code":    "UPDATE_PRODUCT_TAX_INPROGRESS",
				"error_message": "update product tax still in progress",
			},
		}

		if whiteListCompanyType[strconv.Itoa(int(result.Data.CompanyDetail.Items.Type))] {
			if _, ok := whiteListProductTaxStatus[strconv.Itoa(int(result.Data.CompanyDetail.Items.UpdateProductTaxStatus))]; ok {
				return sange.SetError(sange.Forbidden, errors.New("company still updated tax"), "company still updated tax")
			}
		}
	}
	return nil
}
