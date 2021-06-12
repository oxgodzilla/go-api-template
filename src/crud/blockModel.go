package crud

import (
	"github.com/geometry-labs/go-service-template/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type BlockRawModel struct {
	db        *gorm.DB
	model     *models.BlockRaw
	writeChan chan *models.BlockRaw
}

var blockRawModelInstance *BlockRawModel
var blockRawModelOnce sync.Once

func GetBlockRawModel() *BlockRawModel {
	blockRawModelOnce.Do(func() {
		blockRawModelInstance = &BlockRawModel{
			db:        GetPostgresConn().conn,
			model:     &models.BlockRaw{},
			writeChan: make(chan *models.BlockRaw, 1),
		}

		err := blockRawModelInstance.Migrate()
		if err != nil {
			zap.S().Error("BlockModel: Unable create postgres table BlockRaws")
		}
	})
	return blockRawModelInstance
}

func NewBlockRawModel(conn *gorm.DB) *BlockRawModel { // Only for testing
	blockRawModelInstance = &BlockRawModel{
		db:        conn,
		model:     &models.BlockRaw{},
		writeChan: make(chan *models.BlockRaw, 1),
	}
	return blockRawModelInstance
}

func (m *BlockRawModel) GetDB() *gorm.DB {
	return m.db
}

func (m *BlockRawModel) GetModel() *models.BlockRaw {
	return m.model
}

func (m *BlockRawModel) GetWriteChan() chan *models.BlockRaw {
	return m.writeChan
}

func (m *BlockRawModel) Migrate() error {
	// Only using BlockRawORM (ORM version of the proto generated struct) to create the TABLE
	err := m.db.AutoMigrate(models.BlockRawORM{}) // Migration and Index creation
	return err
}

func (m *BlockRawModel) Create(block *models.BlockRaw) *gorm.DB {
	tx := m.db.Create(block)
	return tx
}

func (m *BlockRawModel) Update(oldBlock *models.BlockRaw, newBlock *models.BlockRaw, whereClause ...interface{}) *gorm.DB {
	tx := m.db.Model(oldBlock).Where(whereClause[0], whereClause[1:]).Updates(newBlock)
	return tx
}

func (m *BlockRawModel) Delete(conds ...interface{}) *gorm.DB {
	tx := m.db.Delete(m.model, conds...)
	return tx
}

func (m *BlockRawModel) FindOne(conds ...interface{}) (*models.BlockRaw, *gorm.DB) {
	block := &models.BlockRaw{}
	tx := m.db.Find(block, conds...)
	return block, tx
}

func (m *BlockRawModel) FindAll(conds ...interface{}) (*[]models.BlockRaw, *gorm.DB) {
	blocks := &[]models.BlockRaw{}
	tx := m.db.Scopes().Find(blocks, conds...)
	return blocks, tx
}
