# AP syntax

## Comments

```
//  Single line comment
/*  Multi
    Line
    Comment */
```

## Expressions

All basic expressions should end in a semi colon

## Variables

Variables may not be initialised without a value. Variables cannot be made of anything other than letters

## Assignment

To assign a value to a variable, it should be of the same type, it will not explicitly thow an error in the compiler but will not give correct results or may throw a fatal error

Floats should always have a .0 or another value to the right of the whole number, while integers cannot, even 1.0 would be invalid for integers

```
int value = 1;
float othervalue = 1.0;
bool alsoavalue = true;
```

## Arithmetic operators

Floats cannot be mixed with integers in math

```
int value = 1 + 4 * 3 / 4 - 4;
float othervalue = 1.0 + 4.0 * 3.0 / 4.0 - 4.0;
```

## Bitwise operators

Bitwise operators may be used on any type, if done between two truth values it will give the result of that truth value with the operator

```
int value = 1 & 2; // 1
bool othervalue = true | false; // true
```

## Logical operations

```
a == b
a != b
a > b
a < b
a >= b
a <= b
```

These operations can then be stored in boolean variables

```
bool value = a == b;
bool othervalue = a < b;
```

## Selection

If statements will have its truth condition directly next to it

```
if true {
    // Do stuff
}
```

A variable or expression can also be placed instead

```
if a {
    // Do stuff
}
if false == a {
    // Do stuff
}
```

Bitwise operators can be used also

```
if a | b {
    // Do stuff
}
if a > b & b > c {
    // Do stuff
}
```

## Functions

### Main function

Each program must have a main function, this function cannot be referenced

```
func main() () {
    // Do stuff
}
```

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