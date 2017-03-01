parse
----------

Simple parser constructor inspired by [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) which
describes grammars with composition of go types.

This is work-in-progress project.

## Mocking

``` text
Number       = Either{
    Terminal("0"), Terminal("1"), Terminal("2"),
    Terminal("3"), Terminal("4"), Terminal("5"),
    Terminal("6"), Terminal("7"), Terminal("8"),
    Terminal("9")
}
LeftBracket  = Terminal("(")
RightBracket = Terminal(")")
FixedString  = Chain{Terminal("FixedString"), LeftBracket, Repetition{Number}, RightBracket}
Array        = Chain{Terminal("Array"), LeftBracket, Repetition{Type}, RightBracket}
Type         = Either{Terminal("Date"), Terminal("DateTime"), FixedString, Array}

----------------------------------------------------------------- ->
Actualy go is imperative, we must append() to make type recursive ->
----------------------------------------------------------------- ->

Type  = Either{Terminal("Date"), Terminal("DateTime"), FixedString}
Array = Chain{Terminal("Array"), LeftBracket, Repetition{Type}, RightBracket}
Type = append(Type, Array)
```
