Создание случайных графиков и случайных значений путём перемножения синусов\
Установка:\
```
go get github.com/zongrade/graph
```
Получить случайный график можно с помощью:\
```
graph.CreateRandomGraphicDefault()
```
будет получено,что-то такое:\
![случайный график](https://github.com/zongrade/graph/raw/main/sin_1.png)\
![случайный график](https://github.com/zongrade/graph/raw/main/sin_2.png)\
можно указывать количество количество умножаемых синусов и размер графика в точках
```
graph.CreateRandomGraphic(5, 300)
```
Для получения точек настройте Range в диапазоне 0-100\
```
graph.CreateRandomGraphicDefault().CreateDotsGraphics()
```
будет получено,что-то такое:\
![случайный график](https://github.com/zongrade/graph/raw/main/dots_1.png)\
![случайный график](https://github.com/zongrade/graph/raw/main/dots_2.png)\