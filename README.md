<p align="center">
<a href="https://hackheroes.pl"><img src="https://hackheroes.pl/img/rsz_hackheroes_logo.png"></a>
<img src="internal/assets/icons/logo.png">
</p>

ConstiTutor jest aplikacją, tworzoną do udziału w Ogólnopolskim konkursie
programistycznym [Hack Heroes](https://hackheroes.pl/)

Aplikacja służy do wyszukiwania interesującej użydkownika frazy w Konstytucji
Rzeczypospolitej Polskiej oraz innych aktach prawnych.


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
