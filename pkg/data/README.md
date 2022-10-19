# Pakiet odpowiedzialny za importowanie danych do struktury GO

folder `DATA` zawiera zestaw danych obsługiwanych przez dekoder.

Obsługujemy następujące formaty:

## JSON

pliki z rozszerzeniem `.json` będą konwertowane za pomocą biblioteki
[`encoding/json`](https://pkg.go.dev/encoding/json)

### szablon:

```json
{
        "ActName": "Prawo Dżungli",
        "Rules": [
                {
                        "Identifier": "Prawo Silniejszego",
                        "Text": "Mocniejszy zawsze wygrywa",
                        "Links": [
                                "https://pkg.go.dev"
                        ]
                }
        ]
}
```

## TXT

opracowany przez nasz zespół specjalny format danych dekodowany będzie w następujący
sposób:

deokder będzie iterować po linijkach pliku (pomijając puste linie).
Każda linia powinna zaczynać się jednym z poniższych znaków (w przeciwnym razie może zostać zgłoszony błąd)
- `#` - komentarz (linijka nie jest brana pod uwagę)
- `$` - nazwa aktu prawnego (obecny standard przewiduje nie więcej niż jeden akt prawny na jeden plik txt)
- `&` - nazwa rozdziału
- `@` - nazwa podrozdziału
- `&` - nazwa artykułu
- `*` - tekst artykułu

dekoder będzie wczytywał wszystkie pliki w katalogu DATA oznaczone
rozszerzeniem `.txt`
