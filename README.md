# MyP-Proyecto2
Repositorio para el Proyecto 2 en Modelado y Programación

Lenguaje utilizado: go version go1.11 linux/amd64

Dependencias: [gotk3](https://github.com/gotk3/gotk3)
              [go-mp3](https://github.com/hajimehoshi/go-mp3)
              [oto](https://github.com/hajimehoshi/oto)
              
El proyecto consiste en un reproductor de mp3 con acceso a una base de datos que
incluye la información encontrada en las etiquetas ID3v2.4 en cada mp3. La base
de datos está gestionada con SQLite, e incluye un lenguaje amigable para llevar
a cabo las búsquedas sin tener que utilizar el lenguaje de SQL.

Antes de compilar, gotk3 debe estar instalado en el GOPATH, en la ruta
src/github.com/gotk3/gotk3.   Para instalarlo, basta utilizar el comando

```bash
$ go get github.com/gotk3/gotk3/gtk
```

Para instalarlo en Ubuntu/Debian, se necesitan las siguientes dependencias:
GTK 3.6-3.16, GLib 2.36-2.40, y Cairo 1.10 or 1.12.

Para instrucciones detalladas, consulte: [installation](https://github.com/gotk3/gotk3/wiki#installation)
Las dependencias pueden obtenerse (Ubuntu/Debian) con el comando.

```bash
$ sudo apt-get install libgtk-3-dev libcairo2-dev libglib2.0-dev
```

Para la reproducción de audio se utiliza go-mp3, que a su vez depende de oto,
estos deben de encontrarse en el GOPAH en las rutas
src/github.com/hajimehoshi/go-mp3 y src/github.com/hajimehoshi/oto,
respectivamente.   Una vez en la ruta src/github.com/hajimehoshi/, éstos pueden
obtenerse con los comandos

```bash
$ git clone https://github.com/hajimehoshi/go-mp3.git
$ git clone https://github.com/hajimehoshi/oto.git
```


El paquete oto requiere a su vez libasound2-dev,
que puede obtenerese en Ubunto o Debian con el comando

```bash
$ sudo apt install libasound2-dev
```
