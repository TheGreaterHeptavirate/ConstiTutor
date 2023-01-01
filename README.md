# NOTE 
<image align="left" src="https://user-images.githubusercontent.com/242652/138285004-b27d55b3-163b-4fe3-a8ff-6c34518044bd.png">
Dear Users/visitors
Since the fact that at december 2022 HackHeroes has finished, and the project didn't win anything,
as well as because of lack of interesst in further development,
We've decided to announce its end-of-life.
If you're interessted in resuming the development anyhow, feel free to
contact us for details.

Thanx!
gucio321 <gucio321@users.noreply.github.com> - the Head Developer.
<br clear="all" />


<p align="center">
<a href="https://hackheroes.pl"><img src="https://hackheroes.pl/img/rsz_hackheroes_logo.png"></a>
<img src="internal/assets/icons/logo.png">
</p>

ConstiTutor jest aplikacją, tworzoną do udziału w Ogólnopolskim konkursie
programistycznym [Hack Heroes](https://hackheroes.pl/)

Aplikacja służy do wyszukiwania interesującej użydkownika frazy w Konstytucji
Rzeczypospolitej Polskiej oraz innych aktach prawnych.

## Akty prawne

Obecnie, aplikacja pozwala na wyszukiwanie w następujących aktach prawnych:

- Konstytucja Rzeczypospolitej Polskiej

## Instrukacja instalacji

aby uruchomić program musisz zainstalować kilka komponentów:
1. [golang](https://go.dev)
2. Przejdź do strony frameworku [giu](https://github.com/TheGreaterHeptavirate/giu/tree/constitutor#install)
   i zainstaluj wymagane aplikacje
3. zainstaluj zależności wymagane do uruchomienia [oto](https://github.com/hajimehoshi/oto#prerequisite)
4. pobierz kod źródłowy:
```sh
git clone https://github.com/TheGreaterHeptavirate/constitutor
```
5. W konsoli wejdź do katalogu projektu i zainstaluj zależności
```sh
cd constitutor
go get -d ./...
```

teraz, aby uruchomić program wystarczy wykonać następującą komendę:

```sh
go run cmd/constitutor/main.go
```

Na systemie operacyjnym Linux możesz również spróbować użyć
komendy `make`. Obsługiwane komendy:
- `make setup`
- `make build`
- `make run`
- `make test` - uruchomienie testów jednostkowych (unit testów)
- `make cover` - utworzenie raportu o skuteczności testów jednostkowych
- `make clean`
- `make help`


## Design doc

### TASKLIST

- [X] stworzenie bazy JSONowej dla ustaw (najpierw konstytucji, potem może też do innych ustaw)

[więcej informacji](./pkg/data)

- [X] system wejścia (pkg/core/data)

konwerter jsona do GO

- [X] UI

Użycie frameworku [giu](https://github.com/AllenDang/giu).
Można rozważyć [fyne](https://fyne.io), gdyż jest kompatybilna z androidem.

**UWAGA!** Stosujemy [nieoficialną konwencją struktury projektu GO](https://github.com/golang-standards/project-layout)!
