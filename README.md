# booleang

A structured boolean-logic programming language. It uses a completely different
paradigm from any other language I've seen before, so it might be difficult to
actually use it. The main purpose is for experimenting with boolean logic and
logic circuits, and possibly also for testing the logic behind real life
circuits.

## Example

```
circuit main () -> () {
    # set a to low
    0 -> a;

    # output a
    obit(a);

    # a 1Hz clock
    clock 1s {
        # flip a
        ¬a -> a;
    }
}
```

This example defines a variable `a` which is initally set to `0`.
A variable can either be `1` or `0`.

`a` is outputted using the `obit` builtin function, which means that it
will be persistantly displayed to the terminal. Running the program will produce
this output:

```
a = 0
```

Any other outputted variables would be displayed below.

After outputting it, a 1Hz clock is started. Every second, `a` is flipped,
so `0` will turn into `1` and vice versa. The cool thing about booleang is that
the output changes, so instead of saying `a .... 0`, it will change to say:

```
a = 1
```

Since `a` retains its state, you could think of it as a D-type flip-flop if you
were to recreate the circuit in real life.

## Circuits

You've already seen a circuit in the first example of this document. A circuit
is basically a function, but called a circuit instead to fit in with the
logic circuit theme, and also to emphasise the differences. Like functions,
circuits can take arguments. But, unlike functions in most other languages,
they return values to specific variables too. Look at this example:

```
circuit invert (a) -> (b) {
    ¬a -> b;
}
```

It's a very simple circuit which takes an argument, `a`, and stores the inverse
of it in the output parameter `b`. Think of this circuit like a NOT gate:

![](assets/invert.png)

Circuits can have any number of input or output parameters. In
[booleang.txt](booleang.txt), you can see an example of a
[full adder](https://en.wikipedia.org/wiki/Adder_%28electronics%29#Full_adder)

## Macros

It may seem quite unwieldy for all variables to be a single bit, but in practice
it doesn't matter. Of course, Booleang isn't really a general purpose language,
but either way, you can actually represent any digital data using a number of
bits. For example, the number 9 could be represented as 4 bits. In the example
below, the number 9 is represented using four variables, `b0` being the least
significant bit and `b3` being the most.

```
circuit four () -> () {
    1 -> b0;
    0 -> b1;
    0 -> b2;
    1 -> b3;
}
```

Four lines just to make a number? Seriously?? No - of course not. This is a
perfect use case for _macros_. Macros are used to group together sets of bits
into one name.

```
circuit four () -> () {
    %four (b0, b1, b2, b3);
    (1, 0, 0, 1) -> %four;
}
```

This piece of code does the exact same thing, but in half the space.
