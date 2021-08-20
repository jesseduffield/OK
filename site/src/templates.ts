export const templates = [
  { label: 'Hello, playground', value: `let x = 1;\nlet y = 2;\nputs(x + y);` },
  {
    label: 'Checking Equality',
    value: `let equals = fn(a, b) {
  let x = a >= b;
  let y = lazy b >= a;
  return x && y;
};

puts(equals(1, 2));
puts(equals(3, 3));`,
  },
  {
    label: 'Lazy Evaluation',
    value: `let foo = fn() {
  puts("in foo");
  return false;
}

let bar = fn() {
  puts("in bar");
  return true;
};

let x = lazy foo();
let y = lazy bar();

puts(x && y);`,
  },
  {
    label: 'Operator Precedence',
    value: `puts("Without parentheses:")
puts(5 + 2 * 3);
puts("With parentheses:")
puts(5 + (2 * 3));`,
  },
  {
    label: 'Mapping',
    value: `let doubled = map([1,2,3], fn(e) {
  sleep(1); // sleep one second
  e * 2
});

puts(doubled);

let every = fn(arr, check) {
  let passed = true;
  map(arr, fn(e) {
    switch check(e) { case true: passed = false; } }
  )
  return passed;
};

let result = every([5,2,4,1,3], fn(e) { return e >= 2 });

puts(result);`,
  },
  {
    label: 'Privacy',
    value: `notaclass person {
  pack "I am a stupid piece of shit who should not be doing this"

  field name
  field email
}

let p = new person();

// I acknowledge that I am a stupid piece of shit who should not be doing this
p.name = "Jesse";

puts(p.name);`,
  },
  {
    label: 'Initialising a nac',
    value: `notaclass person {
  field name
  field email

  public init fn(selfish, name, email) {
    selfish.name = name;
    selfish.email = email;
  }

  public tostring fn(selfish) {
    "name: " + selfish.name + ", " + "email: " + selfish.email
  }
}

let p = new person();
p.init("Jesse", "jesse@test.com");
puts(p.tostring());`,
  },

  {
    label: 'Evolution',
    value: `notaclass brgousie {
  public whoami fn(selfish) {
    return "a good-for-nothing aristocrat who likes classes";
  }
}

notaclass person {
  field name
  field email
  field likeclas

  public init fn(selfish, name, email) {
    selfish.name = name;
    selfish.email = email;
    selfish.likeclas = false;
  }

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
p.init("John", "");
puts(p.whoami());
p.makeold();
puts(p.whoami());
`,
  },
];
