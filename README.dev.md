* Добавление нового типа

В proto/decimal/<module>/<version> создать .proto файл с декларацией типа.
Если тип участвует в генезисе, в транзакциях, в результатах запроса - поправить:

*) genesis.proto
 
*) tx.proto (обратить внимание на service Msg)
 
*) query.proto (обратить внимание на service Query)

Выполнить генерацию/регенерацию *.pb.go. Команда требует устновленного докера, настроенного на запуск от обычного пользователя:

`make proto-gen`

!!! ВНИМАНИЕ: эта команда при первом использовании или после очистки образов докера скачивает образ tendermintdev/sdk-proto-gen:v0.7 с hub.docker.com. Если хаб отвалится...

Далее вручную добавляем:

*) в генезисе в дефолтные параметры генезиса, в InitGenesis, создаем функцию валидации и пр. в genesis.go и types/genesis.go

*) всё сохранение-получение и обработка транзакций в keeper: keeper.go, msg_server.go

*) не забыть добавить транзации в handler.go


* Legacy balances (наследство)

Ненулевые балансы адресов, с которых никогда не производились транзакции и поэтому их публичные ключи неизвестны, и поэтому получить их новые адреса невозможно.

Часть модуля coin, в генезисе ключ legacy_balances.

Требует: в генезис модуля bank должен быть включен адрес dx1dw0e0mqxja0q88vm5q9tcxc89hj3vtjltkkw4n - это адрес аккаунта legacy_coin_pool, на котором будет суммированный баланс наследства. После выполнении генезиса этот баланс будет управляться модулем coin, который будет каждому выдавать при запросе.

Программа для получения адресов:

```go
package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	cosmosAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func main() {
	moduleAddress := cosmosAuthTypes.NewModuleAddress("legacy_coin_pool")
	address, err := bech32.ConvertAndEncode("dx", moduleAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("address = %s\n", address)
}
```

Транзакция возврата наследства следующая: кто угодно отправляет транзакцию со следующими параметрами:

*) sender - сам отправитель

*) receiver - новый адрес аккаунта

*) PublicKeyBytes - байты публичного ключа

По байтам публичного ключа получается старый и новый адрес. Если новый адрес совпадает с receiver, а по старому есть баланс - баланс зачисляется на новый адрес.