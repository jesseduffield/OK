# _OK?_

_Try it out on the [playground](https://www.okquestionmark.org/)_

_Watch the [Making-of video](https://www.youtube.com/watch?v=PLGpUsSL0FI&ab_channel=JesseDuffield)_

<p align="center">
  <img src="assets/20220904_154920.jpg" width="700px"/>
  <p align="center"><i>COMPAQ personal Computer BASIC prompt.</i></p>
</p>

## Programming Is Simple Again

<p align="center">
  <img src="assets/quentin.jpg" width="700px"/>
  <p align="center"><i>¿OK?</i></p>
</p>                                


_OK?_ is a modern, dynamically typed programming language with a vision for the future.
_OK?_'s mission is to do away with the needless complexity of today's programming languages and let you focus on what matters: writing code that makes a difference.

### Contents

- [Conditionals](#conditionals)
- [Readable Switches](#readable-switches)
- [Nulls Are Not OK](#nulls-are-not-ok)
- [Error Handling](#error-handling)
- [Readable Logical Operators](#readable-logical-operators)
- [One Comparison Operator](#one-comparison-operator)
- [Dead-simple Operator Precedence](#dead-simple-operator-precedence)
- [Death To Classes](#death-to-classes)
  - [All Fields Are Private](#all-fields-are-private)
  - [No Constructors](#no-constructors)
  - [Evolution Over Composition](#evolution-over-composition)
- [Familiarity Admits Brevity](#familiarity-admits-brevity)
- [Concurrency, Iterated](#concurrency-iterated)
- [Testimonials](#testimonials)
- [How To Get Started](#how-to-get-started)
- [Credits](#credits)

### Conditionals

Very early on in its design, it was decided that _OK?_ would not feature a ternary operator. For those unaware, the ternary operator looks like this:

```js
let a = isprod ? 'prod' : 'dev';
```

Disgusting, we agree. Einstein, Tesla, and Newton all died long ago, so there's really only a handful of humans left on Earth who are capable of parsing that stupifying syntax. What does the question mark mean? What does the colon mean? In _OK?_ we leave those questions for the philosophers and focus on what's important: _writing clean code_.

A language only needs one conditional control flow construct, and in _OK?_, that construct is the switch statement. Switch statements are more versatile and expressive than ternary operators and if statements, and after a while you'll forget those other constructs ever existed. Here's the above statement in idiomatic _OK?_:

```go
let a = switch isprod {
  case true: "prod";
  case false: "dev";
}
```

The switch form, although longer, is unquestionably clearer. Understanding the value of simplicity over complexity is the first step to learning _OK?_

### Readable Switches

Given that switches are so central to _OK?_, we wanted to avoid some common pitfalls around switches found in other languages. In other languages it's common to have a single switch statement take up several pages of an editor with bloated logic being shoved into each switch case, so in _OK?_ we made it that you can only have _one_ statement per case:

If you want to execute multiple statements per switch case, just wrap them in a function:

```go
// INVALID:
switch x {
  case true:
    z = z + 2
  case false:
    x = x + 1;
    y = y - 1; // <-- ERROR: switch blocks can only contain a single statement
}

// VALID:
let onfalse = fn() {
  x = x + 1;
  y = y - 1;
};

switch x {
  case true: z = z + 2;
  case false: onfalse();
}
```

This ensures separation of concerns: the specific per-case logic is factored away, bringing the cases themselves to the forefront.

### Nulls Are Not OK

Null values, famously dubbed the [billion dollar mistake](https://www.infoq.com/presentations/Null-References-The-Billion-Dollar-Mistake-Tony-Hoare/). I want my money back Tony Hoare.

Implicit nullability is a thorn in the side of every self-respecting developer and that's why in _OK?_, null values are _NOT_ OK, _OK?_

We could have omitted null values from _OK?_, but then you might forget how bad they are, so we've included them, represented not as `nil` or `null` but as `NO!`, which is the only sensible answer to the question _Are nulls OK?_

To drive the point home, we've also added a builtin function named `ayok?` which takes a value and returns whether that value is OK.

```go
let x = NO!
puts(ayok?(x)) // prints 'false'

let y = 10
puts(ayok?(y)) // prints 'true'
```

Now when your colleagues ask whether nulls are OK, you'll know exactly what to say.

### Error Handling

In _OK?_, errors are simply values, just like any other value. One value they're particularly similar to is strings, and that's because by convention, they actually are strings. For example:

```go
let divide = fn(a, b) {
  return switch b {
    case 0: [NO!, "cannot divide by zero"];
    default: [a / b, ""];
  };
};

result = divide(5, 0)
switch result[1] {
  case "": puts(result[0])
  default: puts(result[1]) // prints "cannot divide by zero"
}
```

No magic, just arrays and strings.

### Readable Logical Operators

Ever come across a conditional statement that chains a heap of long boolean expressions together?

```go
switch p.isactive() && p.credits() >= reqcreds && p.nottype("Admin") {
  ...
}
```

I think I pulled a neck muscle trying to read that obscenely long line.

In _OK?_, the `&&` and `||` operators can only act on variables, so that it's possible for the reader to understand what's going on.

```go
// ERROR: '&&' operator must act on variables
switch p.isactive() && p.credits() >= reqcreds && p.nottype("Admin") {
  ...
}
```

Here's how the above switch statement would be done in _OK?_

```go
let isactive = p.isactive()
let enoughcr = p.credits() >= reqcreds
let notadmin = p.nottype("Admin")
switch isactive && enoughcr && notadmin {
  ...
}
```

This reflects a central tenet of _OK?_: _be kind to the reader_.

If you need to short-circuit your conditionals, you can use the `lazy` keyword:

```go
let isactive = p.isactive()
let enoughcr = lazy p.credits() >= reqcredits
let notadmin = lazy p.nottype("Admin")
switch isactive && enoughcr && notadmin {
  ...
}
```

If `p.isactive()` returns `true`, then `p.credits()` and `p.usertype()` will never be called.

With this feature you get the best of both worlds: clean, readable code, without sacrificing performance.

### One Comparison Operator

Ever found yourself running in circles looking for the right comparison operator? Instead of wasting time memorising the hieroglyphics of `==`, `!=`, `>`, `<`, `>=`, and `<=`, what if we told you that you only needed one?

_OK?_ features a single comparison operator: `>=`

| In other languages | In _OK?_                                         |
| ------------------ | ------------------------------------------------ |
| `a >= b`           | `a >= b`                                         |
| `a <= b`           | `b >= a`                                        |
| `a > b`            | `!(b >= a)`                                      |
| `a < b`            | `!(a >= b)`                                      |
| `a == b`           | `let x = a >= b; let y = b >= a; x && y`         |
| `a != b`           | `let x = !(a >= b); let y = !(b >= a); x \|\| y` |

This means that instead of juggling six different comparison operators in your head you can focus on what matters: creating great software.

In the _OK?_ standard library we haven't found the need to extract any of these expressions into their own functions, but there's nothing stopping you from rolling your own:

```go
let equals = fn(a, b) {
  let x = a >= b;
  let y = lazy b >= a;
  return x && y;
}
```

### Dead-simple Operator Precedence

in _OK?_, `5 + 2 * 3` evaluates to 21, not 11, because addition and multiplication have equal operator precedence. If you want to evaluate your expression in some other order, you simply need to use parentheses: `5 + (2 * 3)`.

This simple left-to-right default spares you from scrounging around the internet looking for an operator precedence table, and lets you keep your eyes on the code.

### Death To Classes

The authors of _OK?_ watched as object-oriented (OO) languages boomed in popularity, only to find them soon buckling under their own weight. Central to this clinical obesity is the _class_.

A class takes a sensible idea: defining data along with methods that act on that data, and then drives it off a cliff by adding inheritance and subtype polymorphism. It should be no surprise that a bunch of class-obsessed aristocratic oldies in the 60s, who probably spent all their time deciding which child should inherit most of the estate, decided to add a construct named 'class' which revolved around inheritance.

Well, the revolution has finally come! There is no inheritance in _OK?_, in solidarity with all those who fought against the bourgeoisie in years past. Some other languages opted for _structs_ over classes yet despite their developers honourably denouncing OO, old habits die hard with some accidentally uttering 'class' when they mean 'struct'. Some developers even say 'class' deliberately, a dogwhistle to return to the old days when the OO aristocracy still held the developer profession by the throat.

To remove any ambiguity and to ensure full commitment to the death of classes, we've decided to use our own terminology: 'notaclass' , or 'nac' for short.

```go
notaclass person {
  field name
  field email
}
```

Let's deep dive into what makes our nacs special:

#### All Fields Are Private

You don't check how your friend is feeling by prying them open with a crowbar and perusing through their entrails; you just _ask_ them. The same is true in programming. There is no way to mark a field as public in _OK?_ because the public API of a nac should describe _behaviour_, not _state_.

```go
let p = new person()
p.name = "Jesse" // <-- ERROR: access of private field 'name'
```

For extenuating circumstances, you can define a _privacy acknowledgement_ with the `pack` keyword, allowing external code to access a nac's fields if they include the acknowledgement in a comment, preceded by 'I acknowledge that':

```go
notaclass person {
  pack "I am a stupid piece of shit who should not be doing this"

  field name
  field email
}

let p = new person()
// I acknowledge that I am a stupid piece of shit who should not be doing this
p.name = "Jesse"  // <-- No error

```

This makes it easy to find privacy violations with `CTRL+F` and lets you communicate your tolerance level explicitly.

#### No Constructors

What part of `notaclass` don't you understand? Constructors are a class-based thing, and _OK?_ does not have classes. The word `Constructor` also contains the word `struct`, and _OK?_ does not have structs. If you want to define the initial state of a nac, just add it as a separate method:

```go
notaclass person {
  field name
  field email

  public init fn(selfish, name, email) {
    selfish.name = name;
    selfish.email = email;
  }
}

let p = new person()
p.init("Jesse", "jesse@test.com")
```

Notice the first argument in that method: we considered using `self`, `this`, or `me`, for the receiver argument, but felt like these all had connotations that would confuse people if carried over into _OK?_. In _OK?_, receivers are just regular function arguments with no special scoping and no special treatment. But they are still kind of similar to receivers in other languages so we settled on _self-ish_, a sensible middle-ground. It's a word that accurately describes you, if you're the kind of person who disagrees with this convention.

#### Evolution Over Composition

You may be familiar with the phrase 'Composition Over Inheritance'. That was cute but the world has moved on. Composition fixes the _fragile base-class_ problem introduced by Inheritance, only to introduce its own _useless base-component_ problem where you can't shake the feeling that you've actually handicapped yourself a bit by depending on composition for code reuse, especially when the component has no ability to interact with its parent.

It's time we take the next logical step: _Evolution Over Composition_. Instead of thinking in terms of _is-a_ or a _has-a_ relationships, think in terms of _becomes-a_ relationships. How does this work?

To enable evolution, you simply need to define an `evolve` method in your nac. The `evolve` method is invoked after any other of the nac's methods are executed, and it determines if the preconditions have been met to evolve the nac instance into a new nac type.

```go
notaclass brgousie {
  public whoami fn(selfish) {
    return "a good-for-nothing aristocrat who likes classes"
  }
}

notaclass person {
  field name
  field email
  field likeclas

  ...

  public whoami fn(selfish) {
    return selfish.name;
  }

  public makeold fn(selfish) {
    selfish.likeclas = true;
  }

  evolve fn(selfish) {
    switch selfish.likeclas {
      case true:
        return new brgousie();
      default:
        return NO!;
    }
  }
}

let p = new person();
p.init("John", "")
puts(p.whoami()); // prints "John"
p.makeold(); // evolve() method is called behind the scenes
puts(p.whoami()); // prints "a good-for-nothing aristocrat who likes classes"
```

This simple yet powerful feature enables a vast array of possibilities, without the frustration evoked by its predecessors.

### Familiarity Admits Brevity

At some point, the High Counsel Of Programming Conventions got together and decided that variable names need to stretch for miles. It's time to reverse that decision. _Familiarity Admits Brevity_, which is why these days I don't even say goodbye before hanging up on my wife. You should be intimately familiar with your codebase, meaning all of your variables and method names should be short and sweet. You shouldn't need to use juvenile word separators like underscores or camelCase because if you can't capture the meaning of a variable in a single word, that's a sign that you need to refactor.

For this reason, it's idiomatic _OK?_ to limit all variable and method names to eight characters, all in lowercase, and without underscores.

Some example abbreviations:

| invalid                                | valid    |
| -------------------------------------- | -------- |
| characters                             | chars    |
| MaximumPhysicalAddress                 | mxphsadr |
| accrueYesterdaysYield()                | ayy()    |
| hostEnterpriseYellowBorderBackground() | heybbg() |

You may disagree with this idiom, and that's okay, because it's enforced by the compiler. You're welcome.

### Concurrency, Iterated

We've made it this far without talking about iteration or concurrency, and that's because in _OK?_, the two concepts are one and the same.

No modern language can be taken seriously without first-class concurrency support, yet we've found that for all that support, developers still forget to leverage concurrency in the one place that screams for it the most: loops.

#### Introducing _OK?_'s `map` function

_OK?_'s only concurrency construct is also its only loop construct: the `map` function.

`map` behaves similar to map functions in other languages, except for one thing: the callbacks are _always_ executed concurrently.

For example, the following code snippet returns a result in one second:

```go
let doubled = map([1,2,3], fn(e) {
  sleep(1); // sleep one second
  e * 2
}); // [2,4,6] obtained in one second
```

If you just want to compute two expensive values concurrently, you can use `map` with a switch:

```go
let result = map([0,1], fn(e, i) {
  switch i {
  case 0:
    return expnsiv1();
  case 1:
    return expnsiv2();
  }
});
let value1 = result[0];
let value2 = result[1];
```

With the union of concurrency and iteration, the sky is the limit.

```go
let every = fn(arr, check) {
  let passed = true;
  map(arr, fn(e) {
    switch check(e) { case true: passed = false; } }
  )
  return passed;
};

result = every([5,2,4,1,3], fn(e) { return e >= 2 }); // false
```

With this speed, your program's going to finish before you've even started writing it.

### Testimonials

Dave says:

> _OK?_ has transformed me from an angry developer to a happy developer. I used to think composition over inheritance, having my own projects composed of various different languages, never thinking about the power of evolution. Now I've evolved to only use one language: _OK?_

Joel says:

> I used to find _OK?_'s opinionated syntax constraints coercive. Now I find them liberating. Thinking hard about how to fit a complex variable name into eight characters forces me to write code that future me can easily maintain

Sarah says:

> I used to think there were all these features I needed to write quality software, but after a while working with _OK?_ it just clicked: when you stick to the basics, the resulting code is clear and easy to understand.

Jack says:

> If I'm on the _OK?_ subreddit and I see you use the word 'class' when you meant to say 'nac', I'm going to start screaming. And that screaming will not stop until I've found you. By then, you'll be the one screaming.

### How To Get Started

To play around with the language in your browser, you can go to the [playground](https://www.okquestionmark.org/)

To use the language locally, follow the following steps:

1. `git clone` the repo.
2. within the `ok` directory run `go install`.
3. Run `ok` without any arguments to bring up the REPL, or you can run an _OK?_ file with `ok test.ok`.

Happy OK'ing!

### Credits

- Thanks to https://interpreterbook.com/ for helping us create the best new language since assembly
- Thanks to @markbergin77 for the idea about nulls not being OK
- Thanks to @matomatical for showing that even two comparison operators was one too many.
- Thanks to Rob for vectorising Quentyn
