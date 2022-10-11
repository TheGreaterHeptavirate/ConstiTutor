## Design doc

### INTRO

_Obecnie troche słabo ze mną więć bardzo po krótce_

Consti-Tutor - wyszukiwarka po ustawach prawnych Prawa Polskiego

### TASKLIST

- [ ] stworzenie bazy JSONowej dla ustaw (najpierw konstytucji, potem może też do innych ustaw)

JSON powinien wyglądać jakoś tak:

```json
{
        Act: "name of act, e.g. Konstytucja RP",
        Paragraphs: [
                {
                        Index: "1", /* idk czy w jsonie inty są w cudzysłowiu? @garnn?
                        Text: "text of the paragraph",
                        Links: ["links", "to", "explainations", "e.g. youtube"],
                }
                ...
        ],
}
```

to może być w pythonie zrobione, tylko żeby miało jeden format

- [ ] system inputuów (pkg/core/input)

konwerter jsona do GO

- [ ] Backend - (zgodnie z [nieoficialną konwencją GO](https://github.com/golang-standards/project-layout): folder pkg/), polecam `pkg/core`

główny system wyszukiwania e.t.c.

- [ ] UI

może być giu, ale można przetestować coś nowego i użyć [fyne](https://github.com/fyne-io/fyne)
wtedy będziemy mogli łątwo portować na androida (łatwo = `GOOS="android" GOARCH="arm" go build`)
w zależności od czasu

(trzeba zrobić design UI @sirthunderek @TomaszDyrka)
