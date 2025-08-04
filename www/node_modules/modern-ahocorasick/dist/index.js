// src/index.ts
var AhoCorasick = class {
  constructor(keywords) {
    const { failure, gotoFn, output } = this._buildTables(keywords);
    this.gotoFn = gotoFn;
    this.output = output;
    this.failure = failure;
  }
  _buildTables(keywords) {
    const gotoFn = {
      0: {}
    };
    const output = {};
    let state = 0;
    for (const word of keywords) {
      let curr = 0;
      for (const l of word) {
        if (gotoFn[curr] && l in gotoFn[curr]) {
          curr = gotoFn[curr][l];
        } else {
          state++;
          gotoFn[curr][l] = state;
          gotoFn[state] = {};
          curr = state;
          output[state] = [];
        }
      }
      output[curr].push(word);
    }
    const failure = {};
    const xs = [];
    for (const l in gotoFn[0]) {
      const state2 = gotoFn[0][l];
      failure[state2] = 0;
      xs.push(state2);
    }
    while (xs.length > 0) {
      const r = xs.shift();
      if (r !== void 0) {
        for (const l in gotoFn[r]) {
          const s = gotoFn[r][l];
          xs.push(s);
          let state2 = failure[r];
          while (state2 > 0 && !(l in gotoFn[state2])) {
            state2 = failure[state2];
          }
          if (l in gotoFn[state2]) {
            const fs = gotoFn[state2][l];
            failure[s] = fs;
            output[s] = [...output[s], ...output[fs]];
          } else {
            failure[s] = 0;
          }
        }
      }
    }
    return {
      gotoFn,
      output,
      failure
    };
  }
  search(str) {
    let state = 0;
    const results = [];
    for (let i = 0; i < str.length; i++) {
      const l = str[i];
      while (state > 0 && !(l in this.gotoFn[state])) {
        state = this.failure[state];
      }
      if (!(l in this.gotoFn[state])) {
        continue;
      }
      state = this.gotoFn[state][l];
      if (this.output[state].length > 0) {
        const foundStrs = this.output[state];
        results.push([i, foundStrs]);
      }
    }
    return results;
  }
  match(str) {
    let state = 0;
    for (let i = 0; i < str.length; i++) {
      const l = str[i];
      while (state > 0 && !(l in this.gotoFn[state])) {
        state = this.failure[state];
      }
      if (!(l in this.gotoFn[state])) {
        continue;
      }
      state = this.gotoFn[state][l];
      if (this.output[state].length > 0) {
        return true;
      }
    }
    return false;
  }
};
export {
  AhoCorasick as default
};
