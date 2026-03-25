package items

import (
	"net/http"
	"testing"

	apiRunner "QA-TA-SPRING-2026/Task2.1/internal/helpers/api-runner"
	itemsModels "QA-TA-SPRING-2026/Task2.1/internal/managers/items/models"
)

func HttpPostItem(t testing.TB, request itemsModels.CreateItemRequest) *http.Response {
	return apiRunner.GetRunner().Create().Post("/api/1/item").
		JSON(request).
		Expect(t).
		End().Response
}

func HttpGetItemByID(t testing.TB, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get("/api/1/item/" + id).
		Expect(t).
		End().Response
}

func HttpGetItemsBySellerID(t testing.TB, sellerID string) *http.Response {
	return apiRunner.GetRunner().Create().Get("/api/1/" + sellerID + "/item").
		Expect(t).
		End().Response
}

func HttpGetStatistic(t testing.TB, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get("/api/1/statistic/" + id).
		Expect(t).
		End().Response
}

func HttpPostItemWithoutPrice(t testing.TB, request itemsModels.CreateItemRequestWithoutPrice) *http.Response {
	return apiRunner.GetRunner().Create().Post("/api/1/item").
		JSON(request).
		Expect(t).
		End().Response
}
