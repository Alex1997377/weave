# 23.02.2025

## Шаг 2: Описание структур (Data Models) ✅
- Начни с файла internal/core/block.go. Опиши две структуры, о которых мы говорили:
 Header (Index, Timestamp, PrevHash, MerkleRoot, Nonce, Difficulty). ✅
 Block (включает Header, список транзакций и свой Hash). ✅
## Шаг 3: Реализация сериализации и хеширования ✅
- В том же файле напиши методы:
 Serialize() для заголовка (превращаем поля в []byte). Начни с простого bytes.Join, если gob кажется сложным. ✅
 CalculateHash(), который вызывает Serialize и прогоняет результат через sha256.Sum256. ✅
## Шаг 4: Конструктор блока ✅
- Напиши функцию NewBlock(...). Она должна:
 Принять данные (транзакции) и хеш предыдущего блока. ✅
 Создать Header. ✅
- Важно: Вызвать функцию для расчета MerkleRoot (пока можешь сделать заглушку, которая просто склеивает ID транзакций). ✅
## Шаг 5: Цепочка (Blockchain)
- Создай файл internal/core/chain.go. ✅
- Опиши структуру Blockchain (слайс блоков). ✅
- Напиши метод AddBlock(transactions). Он должен находить последний хеш в цепи и передавать его в NewBlock. ✅
- Не забудь про Genesis Block (самый первый блок с пустым PrevHash). ✅
## Шаг 6: Проверка (Validation)
- Напиши метод IsValid(). Это твой главный «тестер»: ✅
- Пройди циклом по цепи. ✅
- Проверь: Current.PrevHash == Previous.Hash. ✅
- Проверь: Current.Hash == Current.CalculateHash(). ✅


## Дополнительно
- Замените типы всех хешей в Block, Header и Transaction на []byte. ✅
- Создайте пакет utils с функцией BytesToHex. ✅
- В интерфейсе Transaction измените GetID() string на GetID() []byte ✅