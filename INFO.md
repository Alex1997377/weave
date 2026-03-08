blockchain-go/
├── cmd/
│   └── blockchain/
│       └── main.go           # Точка входа: инициализация и запуск узла
├── internal/
│   ├── core/                 # Упражнение 1 и 2: Ядро системы
│   │   ├── block.go          # Структуры Block, Header и методы (CalculateHash, Size)
│   │   ├── blockchain.go     # Логика цепочки (AddBlock, IsValid, NewBlockchain)
│   │   ├── transaction.go    # Интерфейс Transaction и базовые методы
│   │   ├── tx_bank.go        # Реализация банковской транзакции
│   │   ├── tx_coinbase.go    # Транзакция вознаграждения майнеру
│   │   ├── tx_pool.go        # Упражнение 2.3: Пул транзакций (Mempool)
│   │   └── errors.go         # Твои кастомные ошибки (BlockchainError)
│   ├── crypto/               # Упражнение 1.2 и 2.2: Криптография
│   │   ├── merkle.go         # Алгоритм дерева Меркла
│   │   ├── proof.go          # Алгоритм майнинга (Proof of Work)
│   │   └── signature.go      # Подписи (ECDSA Sign/Verify)
│   ├── wallet/               # Упражнение 3: Кошельки
│   │   ├── wallet.go         # Генерация ключей и адреса
│   │   └── manager.go        # Хранение и расчет баланса
│   └── store/                # Упражнение 4: Персистентность
│       ├── database.go       # Инициализация (например, BadgerDB или BoltDB)
│       └── repository.go     # Методы SaveBlock, GetBlock, SaveLastHash
├── pkg/
│   └── utils/                # Вспомогательные утилиты
│       ├── convert.go        # BytesToHex, HexToBytes
│       └── serializer.go     # Обертки для gob или json сериализации
├── tests/                    # Упражнение 5: Интеграционные тесты
│   ├── integration_test.go   # Сценарии "от кошелька до блока"
│   └── performance_test.go   # Бенчмарки скорости майнинга
├── go.mod                    # Описание модулей
└── README.md


// Для майнинга полезно иметь быструю проверку с возвратом количества нулей
```func (h Hash) LeadingZeros() int {
    zeros := 0
    
    for _, b := range h {
        if b == 0 {
            zeros += 8
            continue
        }
        
        // Считаем ведущие нули в байте
        for i := 7; i >= 0; i-- {
            if b&(1<<uint(i)) == 0 {
                zeros++
            } else {
                break
            }
        }
        break
    }
    
    return zeros
}```