# AP syntax

## Expresions

All basic expressions should end in a semi colon

## Variables

Variables may not be initialised without a value. Variables cannot be made of anything other than letters

## Assignment

To assign a value to a variable, it should be of the same type, it will not explicitly thow an error in the compiler but will not give correct results or may throw a fatal error

Floats should always have a .0 or another value to the right of the whole number, while integers cannot, even 1.0 would be invalid for integers

```
int value = 1;
float othervalue = 1.0;
```

## Arithmetic operators

Floats cannot be mixed with integers in math

```
int value = 1 + 4 * 3 / 4 - 4;
float othervalue = 1.0 + 4.0 * 3.0 / 4.0 - 4.0;
```

## Functions

### Creating a function

Functions must have brackets for both parameters and a return value no matter if there are none, parameters must have both a type and name, return values must just be given a type

```
func function(int parameter) (int) {
    return parameter;
}
func otherfunction() () {
    int a = 1;
}
```

### Calling a function

```
function();
```

### Passing parameter to a function

```
function() <<< parameter <<< otherparameter;
```