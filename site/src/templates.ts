export const templates = [
  { label: 'Hello, playground', value: `let x = 1;\nlet y = 2;\nputs(x + y);` },
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
