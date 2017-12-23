/*
Package log defines log interface.

Интерфейс logger.Entry используется для разделения библиотеки журналирования (например, logrus)
и кода, который это журналирование использует (например, lib/grpcapi, lib/boltdb).

Интерфейс содержит сигнатуры стандартных методов журналирования и WithField, который у базовой библиотеки (logrus) возвращает внутренний тип, а не этот интерфейс. Поэтому для WithField нужна обертка (см lib/logger).
*/
package log

// Logger is an interface which allows mocks
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
