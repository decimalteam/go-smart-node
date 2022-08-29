# Что необходимо для понимания кода ноды Decimal smart chain

1) Понять консенсус tendermint

2) dsc основан на ethermint, cosmos sdk, tendermint; ethermint основан на evm из go-ethereum и cosmos; cosmos основан на tendermint

3) ethermint используем как готовую evm (go-ethereum) в виде модуля cosmos sdk

4) tendermint напрямую не вызываем, только пользуемся как RPC для запроса информации и отправки сериализованной транзакции

5) активно используем cosmos, нужно обратить внимание на:
- декларацию сообщений и сервисов в proto файлах: на них стоятся типы и транзакции
- типы sdk.Int, sdk.Coin+sdk.Coins, sdk.Context
- модули x/auth, x/bank, почитать про структуру наших модулей ниже
- варианты взаимодействия с нодой: cli, rest, grpc
- cli сделан модульно, конкретные команды и использование флагов прописаны в модулях; есть наши модули, есть модули cosmos
- пример работы с grpc можно увидеть в sdk
- rest признан устаревшим, сейчас rest по порту 1317 - это надстройка над GRPC, ендпоинты можно увидеть в query.proto модулей
- проверка транзакций осуществляется в AnteHandler; есть наши, но основная часть - cosmos, ethermint

# Подключение своей ноды к devnet

0) Код ноды должна быть той же версии (коммита), что и devnet, либо без измненений в логике приложений

1) Собрать ноду через `make install`

2) Выполнить команду `dscd init <node moniker> --chain-id <chain id from genesis>`

3) Скопировать с devnet и положить в `$HOME/.decimal/daemon/config` файлы genesis.json, addrbook.json

4) Запустить ноду `dscd start`, пойдет синхронизация

5) После синхронизации можно запрашивать данные с ноды, отправлять транзакции

# Добавление нового модуля

В proto/decimal/&lt;module&gt;/&lt;version&gt; создать .proto файлы с декларациями (минимум):

*) genesis.proto

*) query.proto

*) tx.proto

*) event.proto

*) различные типы в *.proto

Примерный заголовок
```
syntax = "proto3";
package decimal.*.v1;

import "gogoproto/gogo.proto";

// если необходимы декларации из других файлов
import "decimal/*/v1/*.proto";

option go_package = "bitbucket.org/decimalteam/go-smart-node/x/*/types";

message GenesisState {}
```

Структура модуля. Дерево от x/&lt;module&gt;.

```
.
+-- module.go (AppModule + AppModuleBasic: конструктор, регистрация rest/grpc роутов, cli комманд; это все потом уходит в app.go)
+-- genesis.go (InitGenesis + ExportGenesis)
+-- handler.go (NewHandler: обработка известных типов транзакций)
+-- client
|   +-- cli
|   |   +-- query.go (запросы через cli)
|   |   +-- tx.go (отправка транзакций через cli)
|   +-- rest
|   |   +-- rest.go (регистрация роутов)
|   |   +-- query.go (регистрация хэндлеров gorilla/mux)
+-- keeper
|   +-- keeper.go (операции со стором key-value: загрузка, сохранение, итерация)
|   +-- msg_server.go (обработка транзакций - имплементация описаний service Msg в tx.proto)
|   +-- querier.go (конструктор NewQuerier для module.go)
|   +-- grpc_query.go (обработка запросов - имплементация описаний service Query в query.proto)
+-- simulation
|   +-- simap.go (нужен для тестирования)
+-- types
|   +-- codec.go (регистрация типов транзакций в кодеках)
|   +-- config.go (константы модуля)
|   +-- errors.go (типизированные ошибки с кодами, codespace для облегчения разбора ошибок внешними сервисами, перевода на другие языки)
|   +-- events.go (константы полей при отправке событий)
|   +-- genesis.go (DefaultGenesis(), Validate())
|   +-- keys.go (ключи keeper, используются при сохранении в сторе)
|   +-- msg.go (конструкторы сообщений-транзакций + имплементация методов для удовлетворения типу sdk.Msg: Route, Type, GetSignBytes, GetSigners, ValidateBasic)
|   +-- params.go (нужен при наличии глобальных параметров модуля, которые задаются в генезисе и должны храниться в сторе)

```

# Добавление нового типа и транзакции

В proto/decimal/&lt;module&gt;/&lt;version&gt; создать .proto файл с декларацией типа.
Если тип участвует в генезисе, в транзакциях, в результатах запроса - поправить:

- genesis.proto
- tx.proto (обратить внимание на service Msg)
- query.proto (обратить внимание на service Query)

### Процедура генерации типов до cosmos 0.46
Выполнить генерацию/регенерацию *.pb.go. Команда требует устновленного докера, настроенного на запуск от обычного пользователя:

`make proto-gen`

!!! ВНИМАНИЕ: эта команда при первом использовании или после очистки образов докера скачивает образ tendermintdev/sdk-proto-gen:v0.7 с hub.docker.com. Если хаб отвалится...

### Процедура генерации типов с cosmos 0.46

1) Установить buf, protobuf v3 и утилиты для генерации. Команды описаны в комментарии в `./scripts/protocgen.sh`

2) Для доступа к https://buf.build/ потребуется VPN

3) Запустить `./scripts/protocgen.sh`

### Далее вручную добавляем:

- в генезисе в дефолтные параметры генезиса, в InitGenesis, создаем функцию валидации и пр. в genesis.go и types/genesis.go
- всё сохранение-получение и обработка транзакций в keeper: keeper.go, msg_server.go
- добавить транзации в handler.go

# Запуск ноды как единственной в локальном блокчейне

Должны быть установлены go, jq

Скрипт `init.sh` компилирует бинарник `dscd`. Текущие проблемы в init.sh

- установка claims, которые вызывают проблемы
- установка total supply требует чтобы на сумма монет на всех аккаунтах была равна total supply, поэтому в тестовых целях его лучше убрать

Генезис, конфиги, БД с состояние блокчейна пишутся в `$HOME/.decimal/daemon/`

Для включения REST сервера см. `app.toml`, секция `[api]`

# Legacy balances (наследство)

Ненулевые балансы адресов, с которых никогда не производились транзакции и поэтому их публичные ключи неизвестны, и поэтому получить их новые адреса невозможно.

Модуля legacy, в генезисе ключ legacy/legacy_records. В каждой записи есть информация по балансам, мультисигнатурным кошелькам, NFT.

Требует: в генезис модуля bank должен быть включен адрес dx1dw0e0mqxja0q88vm5q9tcxc89hj3vtjltkkw4n - это адрес аккаунта legacy_coin_pool, на котором будет суммированный баланс наследства. После выполнении генезиса этот баланс будет управляться модулем legacy, который будет при каждой транзакции по публичному ключу определять старый адрес и обрабатывать legacy_records.

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

- sender - сам отправитель
- receiver - новый адрес аккаунта
- PublicKeyBytes - байты публичного ключа

По байтам публичного ключа получается старый и новый адрес. Если новый адрес совпадает с receiver, а по старому есть баланс - баланс зачисляется на новый адрес.