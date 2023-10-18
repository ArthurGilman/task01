Я предполагаю, что у вас postgresql запущен в контейнере Docker и создана пустая таблица

У меня нет информации как у вас запущена кафка и redis , поэтому можете закомменить в main инициализацию cache и kafka, так же kafka.Startreading() . Сервис запустится

export DSN="host=Имя_контейнера_postgres port=5432 user=postgres password="password" dbname=postgres sslmode=disable"

export TABLE = "mytable"
...

docker build -t task

docker run --name task01 -e DSN=$DSN -e TABLE=$TABLE -p 8080:8080 --network net1 task
