package router

import (
	"github.com/geoo115/property-manager/api/accounting"
	"github.com/gin-gonic/gin"
)

func AccountingRouter(rg *gin.RouterGroup) {
	rg.GET("/invoices", accounting.GetInvoices)
	rg.GET("/invoices/:id", accounting.GetInvoiceByID)
	rg.POST("/invoices", accounting.CreateInvoice)
	rg.PUT("/invoices/:id", accounting.UpdateInvoice)
	rg.DELETE("/invoices/:id", accounting.DeleteInvoice)

	rg.GET("/expenses", accounting.GetExpenses)
	rg.GET("/expense/:id", accounting.GetExpenseByID)
	rg.POST("/expense", accounting.CreateExpense)
	rg.PUT("/expense/:id", accounting.UpdateExpense)
	rg.DELETE("/expense/:id", accounting.DeleteExpense)

	rg.GET("/tenant/invoices", accounting.GetInvoicesForTenant)
	rg.GET("/landlord/invoices", accounting.GetInvoicesForLandlord)
}
