package items

import (
	"strconv"
	"testing"

	"QA-TA-SPRING-2026/Task2.1/internal/client/http/items"
	httpHelper "QA-TA-SPRING-2026/Task2.1/internal/helpers/http-helper"
	itemsModels "QA-TA-SPRING-2026/Task2.1/internal/managers/items/models"
)

func defaultStats() itemsModels.Statistics {
	return itemsModels.Statistics{Likes: 1, ViewCount: 1, Contacts: 1}
}

func CreateItem(t testing.TB, name string, price int, sellerID int, expectedStatusCode int) string {
	req := itemsModels.CreateItemRequest{Name: name, Price: price, SellerID: sellerID, Statistics: defaultStats()}
	resp := items.HttpPostItem(t, req)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)
	return body
}

func GetItemByID(t testing.TB, id string, expectedStatusCode int) string {
	resp := items.HttpGetItemByID(t, id)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)
	return body
}

func GetItemsBySellerID(t testing.TB, sellerID int, expectedStatusCode int) string {
	resp := items.HttpGetItemsBySellerID(t, strconv.Itoa(sellerID))
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)
	return body
}

func GetStatistic(t testing.TB, id string, expectedStatusCode int) string {
	resp := items.HttpGetStatistic(t, id)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)
	return body
}

func CreateItemWithoutPrice(t testing.TB, name string, sellerID int, expectedStatusCode int) string {
	req := itemsModels.CreateItemRequestWithoutPrice{Name: name, SellerID: sellerID, Statistics: defaultStats()}
	resp := items.HttpPostItemWithoutPrice(t, req)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)
	return body
}
