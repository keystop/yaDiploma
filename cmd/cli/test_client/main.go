package main 

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "encoding/json"
)

type uR struct {
	Log string `json:"login"`
	Pas string `json:"password"`
}

type bOut struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

func testSign(t string, log string, pas string) {
	fmt.Println("Тест:", t)
	fmt.Println("Адрес", "http://localhost:8080/api/user/"+t)
	u := uR{
		Log: log,
		Pas: pas,
	}
	fmt.Println("Данные:", u)
	reqBody, err := json.Marshal(&u)
	if err != nil {
		print(err)
	}
	a := "http://localhost:8080/api/user/" + t
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "application/json", "", "POST", reqBody)
	fmt.Println("\n", "Со сжатием:")
	makeZipPostRequest(a, "application/json", "", "POST", reqBody)
	fmt.Println("Окончание теста")
}

func testNewOrder(num string, key string) {
	fmt.Println("Тест:", "NewOrder")
	fmt.Println("Адрес", "http://localhost:8080/api/user/orders")

	reqBody := []byte(num)
	a := "http://localhost:8080/api/user/orders"
	fmt.Println("Данные:", num, "Ключ:", key)
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "text/plain", key, "POST", reqBody)
	fmt.Println("\n", "Со сжатием:")
	makeZipPostRequest(a, "text/plain", key, "POST", reqBody)
	fmt.Println("Окончание теста")
}

func testUserOrders(key string) {
	fmt.Println("Тест:", "UserOrders")
	fmt.Println("Адрес", "http://localhost:8080/api/user/orders")

	reqBody := []byte("")
	a := "http://localhost:8080/api/user/orders"
	fmt.Println("Данные:", "", "Ключ:", key)
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "text/plain", key, "GET", reqBody)
	fmt.Println("\n", "Со сжатием:")
	makeZipPostRequest(a, "text/plain", key, "GET", reqBody)
	fmt.Println("Окончание теста")
}

func testUserBalance(key string) {
	fmt.Println("Тест:", "testUserBalance")
	fmt.Println("Адрес", "http://localhost:8080/api/user/balance")

	reqBody := []byte("")
	a := "http://localhost:8080/api/user/balance"
	fmt.Println("Данные:", "", "Ключ:", key)
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "text/plain", key, "GET", reqBody)
	fmt.Println("\n", "Со сжатием:")
	makeZipPostRequest(a, "text/plain", key, "GET", reqBody)
	fmt.Println("Окончание теста")
}

func testUserWithdrawals(key string) {
	fmt.Println("Тест:", "testUserWithdrawals")
	fmt.Println("Адрес", "http://localhost:8080/api/user/balance/withdrawals")

	reqBody := []byte("")
	a := "http://localhost:8080/api/user/balance/withdrawals"
	fmt.Println("Данные:", "", "Ключ:", key)
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "text/plain", key, "GET", reqBody)
	fmt.Println("\n", "Со сжатием:")
	makeZipPostRequest(a, "text/plain", key, "GET", reqBody)
	fmt.Println("Окончание теста")
}

func testBalanceWithdraw(key, ord string, sum float32) {
	fmt.Println("Тест:", "testBalanceWithdraw")
	fmt.Println("Адрес", "http://localhost:8080/api/user/balance/withdraw")
	bo := bOut{
		Order: ord,
		Sum:   sum,
	}
	fmt.Println("Данные:", bo)
	reqBody, err := json.Marshal(&bo)
	if err != nil {
		print(err)
	}
	a := "http://localhost:8080/api/user/balance/withdraw"
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "application/json", key, "POST", reqBody)
	fmt.Println("\n", "Со сжатием:")
	makeZipPostRequest(a, "application/json", key, "POST", reqBody)
	fmt.Println("Окончание теста")
}

func testAddOrderToAccrual(ord string) {
	fmt.Println("Тест:", "testAccrual")
	fmt.Println("Адрес", "http://localhost:8082/api/orders")
	fmt.Println("Данные:", "")
	a := "http://localhost:8082/api/orders"
	reqBody := []byte(`
		{
			"order": "` + ord + `",
			"goods": [
				{
				  "description": "IPHONE",
				  "price": 47399.99
				},
				{
				  "description": "SAMSUNG",
				  "price": 14599.5
				}
			  ]
			}
	`)
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "application/json", "", "POST", reqBody)
	// fmt.Println("\n", "Со сжатием:")
	// makeZipPostRequest(a, "application/json", "", "GET", reqBody)
	fmt.Println("Окончание теста")
}

func testAccrual(ord string) {
	fmt.Println("Тест:", "testAccrual")
	fmt.Println("Адрес", "http://localhost:8082/api/orders/"+ord)
	fmt.Println("Данные:", "")
	a := "http://localhost:8082/api/orders/" + ord
	reqBody := []byte("")
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "application/json", "", "GET", reqBody)
	// fmt.Println("\n", "Со сжатием:")
	// makeZipPostRequest(a, "application/json", "", "GET", reqBody)
	fmt.Println("Окончание теста")
}

func testAccrualPost(ord string) {
	fmt.Println("Тест:", "testAccrual")
	fmt.Println("Адрес", "http://localhost:8082/api/orders")
	fmt.Println("Данные:", "")
	a := "http://localhost:8082/api/orders"
	reqBody := []byte(`{"order": "` + ord + `"}`)
	fmt.Println("\n", "Без сжатия:")
	makeRequest(a, "application/json", "", "POST", reqBody)
	// fmt.Println("\n", "Со сжатием:")
	// makeZipPostRequest(a, "application/json", "", "GET", reqBody)
	fmt.Println("Окончание теста")
}

func makeRequest(address, ctype, key, rtype string, b []byte) {
	client := &http.Client{}
	req, _ := http.NewRequest(rtype, address, bytes.NewReader(b))
	req.Header.Add("Content-Type", ctype)
	req.Header.Add("Authorization", key)
	// r, err := http.Post(a, t, bytes.NewBuffer(b)) //bytes.NewBuffer(reqBody))
	r, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		print(err)
	}
	printResult(r.Body, r)
}

func makeZipPostRequest(address, ctype, key, rtype string, reqBody []byte) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	gz.Write(reqBody)
	gz.Flush()
	gz.Close()

	client := &http.Client{}
	req, _ := http.NewRequest(rtype, address, bytes.NewReader(b.Bytes()))
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Type", ctype)
	req.Header.Add("Authorization", key)

	r, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	gzr, err := gzip.NewReader(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	printResult(gzr, r)
}

func printResult(body io.Reader, r *http.Response) {
	text, err := io.ReadAll(body)
	if err != nil {
		print(err)
	}
	defer r.Body.Close()
	fmt.Printf("%s\n", r.Header)
	fmt.Printf("%s\n", text)
	fmt.Printf("%d\n", r.StatusCode)
}

// func makeGetPing() {
// 	client := &http.Client{
// 		CheckRedirect: noRedirect,
// 	}
// 	//req, _ := http.NewRequest("GET", "http://localhost:8080/14afc95e687fa093f0edfa25de0766cd", nil)
// 	req, _ := http.NewRequest("GET", "http://localhost:8080/ping", nil)
// 	// response, err := http.Get("http://localhost:8080/04f51bcd17361670c1dc6d94cbbd0efe")

// 	req.AddCookie(&http.Cookie{
// 		Name:  "UserID",
// 		Value: "dc3b5af8713f9d0c1f2dc708e8b2f038",
// 	})

// 	response, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	defer response.Body.Close()

// 	fmt.Printf("%d\n", response.StatusCode)
// }

func main() {
	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200(на новой базе) 409(на старой), 409", "успешно, уже есть")
	// testSign("register", "Aleha", "123123213")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200(на новой базе) 409(на старой), 409", "успешно, уже есть")
	// testSign("register", "Kartoha", "457457457457")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200, 200", "Успешно")
	// testSign("login", "Aleha", "123123213")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200, 200", "Успешно")
	// testSign("login", "Kartoha", "457457457457")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 401, 401", "Неверная пара логи, пароль")
	// testSign("login", "Karas", "457457457457")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 422, 422", "Не верный формат заказа")
	// testNewOrder("4561261212345464", "Bearer 6756be86f17b6853cb7ae5bb78729977")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 202, 200", "Создан, уже есть")
	// testNewOrder("4561261212345467", "Bearer 6756be86f17b6853cb7ae5bb78729977")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 409, 409", "Загружено другим пользователем")
	// testNewOrder("4561261212345467", "Bearer 599b1f0c421ea16cfa0c8ae7c15d9ec2")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 202, 200", "Создан, уже есть")
	// testNewOrder("4561261212345467", "Bearer 6756be86f17b6853cb7ae5bb78729977")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 401, 401", "Пользователь не авторизован")
	// testUserOrders("Bearer 6b98c42394f9ce2763e152c0b5223")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200, 200", "Заказы пользователя")
	// testUserOrders("Bearer 6756be86f17b6853cb7ae5bb78729977")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200, 200", "Баланс пользователя")
	// testUserBalance("Bearer 6756be86f17b6853cb7ae5bb78729977")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200, 200", "Списания пользователя")
	// testUserWithdrawals("Bearer 6756be86f17b6853cb7ae5bb78729977")
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 402, 402", "Списания пользователя")
	// testBalanceWithdraw("Bearer 6756be86f17b6853cb7ae5bb78729977", "", 200000)
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// fmt.Println("--------------------------------------------------------------------------------------------------------")
	// fmt.Println("Ожидаемый результат 200, 200", "Списания пользователя")
	// testBalanceWithdraw("Bearer 6756be86f17b6853cb7ae5bb78729977", "1230", 2500)
	// fmt.Println("--------------------------------------------------------------------------------------------------------")

	// _________________________________________ТЕСТ СИТСЕМЫ НАЧИСЛЕНИЙ ___________________________________________________

	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200, 200", "Вознаграждение пользователя")
	testAddOrderToAccrual("3021351832")
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200, 200", "Вознаграждение пользователя")
	testAccrual("3021351832")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200, 200", "Вознаграждение пользователя")
	testAddOrderToAccrual("4561261212345467")
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("--------------------------------------------------------------------------------------------------------")
	fmt.Println("Ожидаемый результат 200, 200", "Вознаграждение пользователя")
	testAccrual("4561261212345467")
	fmt.Println("--------------------------------------------------------------------------------------------------------")

	// testAccrualPost("1230")

}
