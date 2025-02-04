package replication

import (
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ReplicatedProduct struct {
	ID           string `json:"id"`
	Nombre       string `json:"nombre"`
	Catidad      int    `json:"cantidad"`
	CodigoBarras string `json:"codigobarras"`
}

var (
	replicatedProducts []ReplicatedProduct
	repMu              sync.Mutex
	lastSync           time.Time
)

func CreateReplicatedProduct(c *gin.Context) {
	var newProd ReplicatedProduct
	if err := c.ShouldBindJSON(&newProd); err != nil {
		c.JSON(400, gin.H{"error": "Error al cargar los datos enviados"})
		return
	}

	repMu.Lock()
	replicatedProducts = append(replicatedProducts, newProd)
	lastSync = time.Now()
	repMu.Unlock()

	c.JSON(201, gin.H{"status": "Producto replicado creado", "product": newProd})
}

func ListReplicatedProducts(c *gin.Context) {
	repMu.Lock()
	defer repMu.Unlock()

	c.JSON(200, gin.H{"products": replicatedProducts})
}

func GetReplicatedProductByID(c *gin.Context) {
	id := c.Param("id")

	repMu.Lock()
	defer repMu.Unlock()

	for _, prod := range replicatedProducts {
		if prod.ID == id {
			c.JSON(200, gin.H{"product": prod})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Producto replicado no encontrado"})
}

func UpdateReplicatedProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProd ReplicatedProduct
	if err := c.ShouldBindJSON(&updatedProd); err != nil {
		c.JSON(400, gin.H{"error": "Error al enivar los datos enviados"})
		return
	}

	repMu.Lock()
	defer repMu.Unlock()

	for i, prod := range replicatedProducts {
		if prod.ID == id {
			replicatedProducts[i] = updatedProd
			lastSync = time.Now()
			c.JSON(200, gin.H{"status": "Producto replicado actualizado correctamente :D", "product": updatedProd})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Producto replicado no encontrado :("})
}

func DeleteReplicatedProduct(c *gin.Context) {
	id := c.Param("id")

	repMu.Lock()
	defer repMu.Unlock()

	for i, prod := range replicatedProducts {
		if prod.ID == id {
			replicatedProducts = append(replicatedProducts[:i], replicatedProducts[i+1:]...)
			lastSync = time.Now()
			c.JSON(200, gin.H{"status": "Producto replicado eliminado"})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Producto replicado no encontrado"})
}

func ShortPollingEndpoint(c *gin.Context) {
	repMu.Lock()
	defer repMu.Unlock()

	c.JSON(200, gin.H{
		"products":  replicatedProducts,
		"last_sync": lastSync,
	})
}

func LongPollingEndpoint(c *gin.Context) {
	repMu.Lock()
	currentSync := lastSync
	repMu.Unlock()

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			repMu.Lock()
			if !lastSync.Equal(currentSync) {
				repMu.Unlock()
				c.JSON(200, gin.H{
					"products":  replicatedProducts,
					"last_sync": lastSync,
				})
				return
			}
			repMu.Unlock()
		case <-timeout:
			c.JSON(200, gin.H{"message": "No se hay cambios"})
			return
		}
	}
}
