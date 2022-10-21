package controller

import (
	"net/http"

	"github.com/Pet002/Project-sa-65/entity"
	"github.com/gin-gonic/gin"
)

//------------------------------- Company -----------------------------------

// POST /company
func CreateCompany(c *gin.Context) {
	var MedicineCompany entity.MedicineCompany
	if err := c.ShouldBindJSON(&MedicineCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := entity.DB().Create(&MedicineCompany).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": MedicineCompany})
}

// GET /company/:id
func GetCompany(c *gin.Context) {
	var MedicineCompany entity.MedicineCompany
	id := c.Param("id")
	if err := entity.DB().Raw("SELECT * FROM medicine_types WHERE id = ?", id).Scan(&MedicineCompany).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": MedicineCompany})
}

// GET /company
func ListCompany(c *gin.Context) {
	var MedicineCompany []entity.MedicineCompany
	if err := entity.DB().Raw("SELECT * FROM medicine_types").Scan(&MedicineCompany).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": MedicineCompany})
}

// DELETE /company/:id
func DeleteCompany(c *gin.Context) {
	id := c.Param("id")
	if tx := entity.DB().Exec("DELETE FROM medicine_types WHERE id = ?", id); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MedicineType not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": id})
}

// PATCH /company
func UpdateCompany(c *gin.Context) {
	var MedicineCompany entity.MedicineCompany
	if err := c.ShouldBindJSON(&MedicineCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("id = ?", MedicineCompany.ID).First(&MedicineCompany); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type not found"})
		return
	}

	if err := entity.DB().Save(&MedicineCompany).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": MedicineCompany})
}

//----------------------------------- order -------------------------------

// POST /order
func CreateOrder(c *gin.Context) {

	var medicine entity.Medicine
	var employee entity.Employee
	var company entity.MedicineCompany
	var order entity.MedicineOrder

	// ผลลัพธ์ที่ได้จากขั้นตอนที่  จะถูก bind เข้าตัวแปร order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 8: ค้นหา employee ด้วย id
	if tx := entity.DB().Where("id = ?", order.EmployeeID).First(&employee); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee not found"})
		return
	}

	// 10: ค้นหา company ด้วย id
	if tx := entity.DB().Where("id = ?", order.MedicineCompanyID).First(&company); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company not found"})
		return
	}

	// 11: ค้นหา medicine ด้วย id
	if tx := entity.DB().Where("id = ?", order.MedicineID).First(&medicine); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "medicine not found"})
		return
	}
	// 12: สร้าง medicine
	or := entity.MedicineOrder{
		Employee: employee,        // โยงความสัมพันธ์กับ Entity Employee
		Medicine:  medicine,         // โยงความสัมพันธ์กับ Entity Medicine
		MedicineCompany:	company,   // โยงความสัมพันธ์กับ Entity Company
		OrderAmount:     order.OrderAmount,   // ตั้งค่าฟิลด์ Amount
		OrderTime:	order.OrderTime,    // ตั้งค่าฟิลด์ ordertime
	}

	// 13: บันทึก
	if err := entity.DB().Create(&or).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": or})
}

// GET /order/:id
func GetOrder(c *gin.Context) {
	var order entity.MedicineOrder
	id := c.Param("id")
	if err := entity.DB().Preload("Employee").Preload("Medicine").Preload("Company").Raw("SELECT * FROM medicines WHERE id = ?", id).Find(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

// GET /order
func ListOrder(c *gin.Context) {
	var order []entity.MedicineOrder
	if err := entity.DB().Preload("Employee").Preload("Medicine").Preload("Company").Raw("SELECT * FROM medicines").Find(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// DELETE /order/:id
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if tx := entity.DB().Exec("DELETE FROM medicines WHERE id = ?", id); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": id})
}

// PATCH /order
func UpdateOrder(c *gin.Context) {
	var order entity.MedicineOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("id = ?", order.ID).First(&order); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
		return
	}

	if err := entity.DB().Save(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}
