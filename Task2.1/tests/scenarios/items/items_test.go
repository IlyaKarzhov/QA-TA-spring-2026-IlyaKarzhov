package items

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"QA-TA-SPRING-2026/Task2.1/internal/managers/items"
	itemsModels "QA-TA-SPRING-2026/Task2.1/internal/managers/items/models"
	"QA-TA-SPRING-2026/Task2.1/internal/utils"
)

type ItemsTestSuite struct {
	suite.Suite
}

func TestItemsSuite(t *testing.T) {
	suite.RunSuite(t, new(ItemsTestSuite))
}

func (s *ItemsTestSuite) BeforeAll(t provider.T) {
	os.Setenv("API_URL", "https://qa-internship.avito.com") //nolint:errcheck,gosec
	utils.LoadEnv()
}

func nonExistentUUID() string {
	return uuid.New().String()
}

// randomSellerID returns a random sellerId in range 111111–999999 using crypto/rand.
func randomSellerID() int {
	const min = 111111
	const max = 999999
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		panic("randomSellerID: crypto/rand error: " + err.Error())
	}
	return int(n.Int64()) + min
}

// extractCreatedID parses UUID from POST /api/1/item status response.
// e.g. "Сохранили объявление - <uuid>"
func extractCreatedID(status string) string {
	const sep = " - "
	idx := strings.LastIndex(status, sep)
	if idx < 0 {
		return ""
	}
	return strings.TrimSpace(status[idx+len(sep):])
}

// ========== Positive tests ==========

func (s *ItemsTestSuite) TestCreateItemPositive(t provider.T) {
	t.Title("Create item with valid data")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Create Item"})
	t.Label(&allure.Label{Name: string(allure.Severity), Value: string(allure.CRITICAL)})

	sellerID := randomSellerID()
	price := 1000
	name := "TestItem_" + utils.RandomString(8)

	var createdID string
	var getResp []itemsModels.ItemResponse

	t.WithNewStep("POST item", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		createdID = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(createdID, "UUID must be present in status")
	})

	t.WithNewStep("GET item by ID", func(sCtx provider.StepCtx) {
		getBody := items.GetItemByID(t, createdID, 200)
		t.WithAttachments(allure.NewAttachment("get_response", allure.MimeType("application/json"), []byte(getBody)))
		err := json.Unmarshal([]byte(getBody), &getResp)
		sCtx.Require().NoError(err)
		sCtx.Require().Len(getResp, 1)
	})

	t.WithNewStep("check fields", func(sCtx provider.StepCtx) {
		sCtx.Require().Equal(name, getResp[0].Name)
		sCtx.Require().Equal(price, getResp[0].Price)
		sCtx.Require().Equal(sellerID, getResp[0].SellerID)
	})
}

func (s *ItemsTestSuite) TestGetItemByIDPositive(t provider.T) {
	t.Title("Get existing item by ID")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Get Item"})

	sellerID := randomSellerID()
	price := 1500
	name := "TestGet_" + utils.RandomString(6)

	var createdID string
	var getResp []itemsModels.ItemResponse

	t.WithNewStep("POST item", func(sCtx provider.StepCtx) {
		createBody := items.CreateItem(t, name, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("create_response", allure.MimeType("application/json"), []byte(createBody)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(createBody), &resp)
		sCtx.Require().NoError(err)
		createdID = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(createdID)
	})

	t.WithNewStep("GET item by ID", func(sCtx provider.StepCtx) {
		getBody := items.GetItemByID(t, createdID, 200)
		t.WithAttachments(allure.NewAttachment("get_response", allure.MimeType("application/json"), []byte(getBody)))
		err := json.Unmarshal([]byte(getBody), &getResp)
		sCtx.Require().NoError(err)
		sCtx.Require().Len(getResp, 1)
	})

	t.WithNewStep("check fields", func(sCtx provider.StepCtx) {
		sCtx.Require().Equal(name, getResp[0].Name)
		sCtx.Require().Equal(price, getResp[0].Price)
		sCtx.Require().Equal(sellerID, getResp[0].SellerID)
	})
}

func (s *ItemsTestSuite) TestGetItemsBySellerIDPositive(t provider.T) {
	t.Title("Get all items by seller ID")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Get Items by Seller"})

	sellerID := randomSellerID()
	name1 := "Item1_" + utils.RandomString(4)
	name2 := "Item2_" + utils.RandomString(4)
	price := 1000

	var id1, id2 string
	var listResp []itemsModels.ItemResponse

	t.WithNewStep("POST item 1", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name1, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("create_response1", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		id1 = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(id1)
	})

	t.WithNewStep("POST item 2", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name2, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("create_response2", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		id2 = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(id2)
	})

	t.WithNewStep("GET items by sellerID", func(sCtx provider.StepCtx) {
		body := items.GetItemsBySellerID(t, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("list_response", allure.MimeType("application/json"), []byte(body)))
		err := json.Unmarshal([]byte(body), &listResp)
		sCtx.Require().NoError(err)
	})

	t.WithNewStep("both items in list", func(sCtx provider.StepCtx) {
		sCtx.Require().GreaterOrEqual(len(listResp), 2)
		found1, found2 := false, false
		for _, it := range listResp {
			if it.ID == id1 {
				found1 = true
				sCtx.Require().Equal(name1, it.Name)
				sCtx.Require().Equal(price, it.Price)
				sCtx.Require().Equal(sellerID, it.SellerID)
			}
			if it.ID == id2 {
				found2 = true
				sCtx.Require().Equal(name2, it.Name)
				sCtx.Require().Equal(price, it.Price)
				sCtx.Require().Equal(sellerID, it.SellerID)
			}
		}
		sCtx.Require().True(found1, "Item 1 not found in list")
		sCtx.Require().True(found2, "Item 2 not found in list")
	})
}

func (s *ItemsTestSuite) TestGetStatisticPositive(t provider.T) {
	t.Title("Get statistics for existing item")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Statistics"})

	sellerID := randomSellerID()
	price := 1000
	name := "StatTest_" + utils.RandomString(6)

	var createdID string
	var statResp []itemsModels.StatisticResponse

	t.WithNewStep("POST item", func(sCtx provider.StepCtx) {
		createBody := items.CreateItem(t, name, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("create_response", allure.MimeType("application/json"), []byte(createBody)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(createBody), &resp)
		sCtx.Require().NoError(err)
		createdID = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(createdID)
	})

	t.WithNewStep("GET statistics", func(sCtx provider.StepCtx) {
		statBody := items.GetStatistic(t, createdID, 200)
		t.WithAttachments(allure.NewAttachment("stat_response", allure.MimeType("application/json"), []byte(statBody)))
		err := json.Unmarshal([]byte(statBody), &statResp)
		sCtx.Require().NoError(err)
		sCtx.Require().Len(statResp, 1)
	})

	t.WithNewStep("stats >= 0", func(sCtx provider.StepCtx) {
		sCtx.Require().GreaterOrEqual(statResp[0].ViewCount, 0)
		sCtx.Require().GreaterOrEqual(statResp[0].Contacts, 0)
		sCtx.Require().GreaterOrEqual(statResp[0].Likes, 0)
	})
}

// ========== Negative tests ==========

func (s *ItemsTestSuite) TestCreateItemWithoutName(t provider.T) {
	t.Title("Create item: name missing, expect 400")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})

	sellerID := randomSellerID()
	price := 1000

	t.WithNewStep("POST name='' (empty)", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, "", price, sellerID, 400)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var errResp itemsModels.ErrorResponse
		err := json.Unmarshal([]byte(body), &errResp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(errResp.Result.Message)
	})
}

func (s *ItemsTestSuite) TestCreateItemWithZeroPrice(t provider.T) {
	// BUG-1: price=0 rejected as missing. See BUGS.md.
	t.Title("Create item: price=0 — BUG-1")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})
	t.Label(&allure.Label{Name: string(allure.Severity), Value: string(allure.NORMAL)})

	sellerID := randomSellerID()
	name := "ZeroPrice_" + utils.RandomString(6)

	t.WithNewStep("POST price=0 (BUG-1: server returns 400)", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name, 0, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})
}

func (s *ItemsTestSuite) TestCreateItemWithoutPriceField(t provider.T) {
	t.Title("Create item: price missing, expect 400")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})

	sellerID := randomSellerID()
	name := "NoPriceField_" + utils.RandomString(6)

	t.WithNewStep("POST without price field", func(sCtx provider.StepCtx) {
		body := items.CreateItemWithoutPrice(t, name, sellerID, 400)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var errResp itemsModels.ErrorResponse
		err := json.Unmarshal([]byte(body), &errResp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(errResp.Result.Message)
	})
}

func (s *ItemsTestSuite) TestCreateItemWithoutSellerID(t provider.T) {
	t.Title("Create item: sellerID=0, expect 400")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})

	name := "NoSeller"
	price := 100

	t.WithNewStep("POST sellerID=0", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name, price, 0, 400)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var errResp itemsModels.ErrorResponse
		err := json.Unmarshal([]byte(body), &errResp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(errResp.Result.Message)
	})
}

func (s *ItemsTestSuite) TestCreateItemWithNegativePrice(t provider.T) {
	// BUG-2: negative price accepted by server. See BUGS.md.
	t.Title("Create item: price=-100, expect 400 — BUG-2")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})

	sellerID := randomSellerID()

	t.WithNewStep("POST price=-100 (BUG-2: server returns 200)", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, "NegPrice", -100, sellerID, 400)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var errResp itemsModels.ErrorResponse
		err := json.Unmarshal([]byte(body), &errResp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(errResp.Result.Message)
	})
}

func (s *ItemsTestSuite) TestCreateItemWithSellerIDOutOfRecommendedRange(t provider.T) {
	// Range 111111-999999 is a recommendation only; server has no such validation.
	t.Title("Create item: sellerID outside recommended range")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Corner Cases"})

	t.WithNewStep("POST sellerID=100000", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, "OutRange", 100, 100000, 200)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})
}

func (s *ItemsTestSuite) TestGetItemByIDNotFound(t provider.T) {
	t.Title("Get item by ID: non-existent, expect 404")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})

	id := nonExistentUUID()

	t.WithNewStep("GET non-existent UUID", func(sCtx provider.StepCtx) {
		body := items.GetItemByID(t, id, 404)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		sCtx.Require().NotEmpty(body)
	})
}

func (s *ItemsTestSuite) TestGetItemsBySellerIDNotFound(t provider.T) {
	t.Title("Get items: no listings, expect empty")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})

	sellerID := randomSellerID()

	t.WithNewStep("GET unused sellerID", func(sCtx provider.StepCtx) {
		body := items.GetItemsBySellerID(t, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var listResp []itemsModels.ItemResponse
		err := json.Unmarshal([]byte(body), &listResp)
		sCtx.Require().NoError(err)
		sCtx.Require().Empty(listResp, "Seller with no listings must return empty array")
	})
}

func (s *ItemsTestSuite) TestGetStatisticNotFound(t provider.T) {
	t.Title("Get stats: non-existent item, expect 404")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Negative Scenarios"})

	id := nonExistentUUID()

	t.WithNewStep("GET non-existent UUID", func(sCtx provider.StepCtx) {
		body := items.GetStatistic(t, id, 404)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		sCtx.Require().NotEmpty(body)
	})
}

// ========== Corner cases ==========

func (s *ItemsTestSuite) TestCreateItemWithBoundaryPrice(t provider.T) {
	t.Title("Price boundary: 1 and max int32")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Corner Cases"})

	sellerID := randomSellerID()
	name := "BoundaryPrice"

	t.WithNewStep("POST price=1", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name, 1, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response_min", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})

	t.WithNewStep("POST price=2147483647", func(sCtx provider.StepCtx) {
		maxPrice := 2147483647
		body := items.CreateItem(t, name+"max", maxPrice, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response_max", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})
}

func (s *ItemsTestSuite) TestCreateItemWithBoundarySellerID(t provider.T) {
	t.Title("SellerID boundary: 111111 and 999999")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Corner Cases"})

	price := 100

	t.WithNewStep("POST sellerID=111111", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, "SellerMin", price, 111111, 200)
		t.WithAttachments(allure.NewAttachment("response_min", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})

	t.WithNewStep("POST sellerID=999999", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, "SellerMax", price, 999999, 200)
		t.WithAttachments(allure.NewAttachment("response_max", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})
}

func (s *ItemsTestSuite) TestIdempotentCreate(t provider.T) {
	t.Title("Idempotency: same POST twice gives different IDs")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Corner Cases"})

	sellerID := randomSellerID()
	price := 100
	name := "Idempotent_" + utils.RandomString(5)

	var id1, id2 string

	t.WithNewStep("First create", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response1", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		id1 = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(id1)
	})

	t.WithNewStep("POST same data again", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response2", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		id2 = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(id2)
	})

	t.WithNewStep("IDs must differ", func(sCtx provider.StepCtx) {
		sCtx.Require().NotEqual(id1, id2)
	})
}

// TC-COR-02
func (s *ItemsTestSuite) TestCreateItemWithBoundaryName(t provider.T) {
	t.Title("Name boundary: 1 char and 255 chars")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Corner Cases"})

	sellerID := randomSellerID()
	price := 100

	t.WithNewStep("POST name of 1 char", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, "A", price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})

	t.WithNewStep("POST name of 255 chars", func(sCtx provider.StepCtx) {
		longName := strings.Repeat("a", 255)
		body := items.CreateItem(t, longName, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})
}

// TC-COR-05
func (s *ItemsTestSuite) TestCreateItemWithSpecialCharsName(t provider.T) {
	t.Title("Special characters in name field")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "Corner Cases"})

	sellerID := randomSellerID()

	t.WithNewStep("POST name='!@#$%^&*()'", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, "!@#$%^&*()", 100, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("response", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		sCtx.Require().NotEmpty(extractCreatedID(resp.Status))
	})
}

// TC-POS-05 — E2E
func (s *ItemsTestSuite) TestE2ECreateGetStatistic(t provider.T) {
	t.Title("E2E: create item, get by ID, get statistics")
	t.Label(&allure.Label{Name: string(allure.Feature), Value: "E2E"})
	t.Label(&allure.Label{Name: string(allure.Severity), Value: string(allure.CRITICAL)})

	sellerID := randomSellerID()
	price := 500
	name := "E2E_" + utils.RandomString(6)

	var createdID string
	var getResp []itemsModels.ItemResponse
	var statResp []itemsModels.StatisticResponse

	t.WithNewStep("POST item", func(sCtx provider.StepCtx) {
		body := items.CreateItem(t, name, price, sellerID, 200)
		t.WithAttachments(allure.NewAttachment("create_response", allure.MimeType("application/json"), []byte(body)))
		var resp itemsModels.CreateItemResponse
		err := json.Unmarshal([]byte(body), &resp)
		sCtx.Require().NoError(err)
		createdID = extractCreatedID(resp.Status)
		sCtx.Require().NotEmpty(createdID)
	})

	t.WithNewStep("GET item by ID", func(sCtx provider.StepCtx) {
		body := items.GetItemByID(t, createdID, 200)
		t.WithAttachments(allure.NewAttachment("get_response", allure.MimeType("application/json"), []byte(body)))
		err := json.Unmarshal([]byte(body), &getResp)
		sCtx.Require().NoError(err)
		sCtx.Require().Len(getResp, 1)
		sCtx.Require().Equal(name, getResp[0].Name)
		sCtx.Require().Equal(price, getResp[0].Price)
		sCtx.Require().Equal(sellerID, getResp[0].SellerID)
		sCtx.Require().Equal(createdID, getResp[0].ID)
	})

	t.WithNewStep("GET statistics", func(sCtx provider.StepCtx) {
		body := items.GetStatistic(t, createdID, 200)
		t.WithAttachments(allure.NewAttachment("stat_response", allure.MimeType("application/json"), []byte(body)))
		err := json.Unmarshal([]byte(body), &statResp)
		sCtx.Require().NoError(err)
		sCtx.Require().Len(statResp, 1)
		sCtx.Require().GreaterOrEqual(statResp[0].ViewCount, 0)
		sCtx.Require().GreaterOrEqual(statResp[0].Contacts, 0)
		sCtx.Require().GreaterOrEqual(statResp[0].Likes, 0)
	})
}
