import { escapeHTML } from '../utils/content.js';

// A small, non-executing interpreter for the directives used by the original
// Go HTML templates. No eval, generated JavaScript, or user-supplied templates.
export class SafeHTML {
  constructor(value) { this.value = value; }
}
export const safeHTML = value => new SafeHTML(value);
const truth = value => Array.isArray(value) ? value.length > 0 : Boolean(value);

function expression(source, context) {
  const tokens = source.match(/"(?:\\.|[^"\\])*"|[()]|[^\s()]+/g) || [];
  let position = 0;
  function group(nested = false) {
    const terms = [];
    while (position < tokens.length) {
      const token = tokens[position++];
      if (token === ')') break;
      terms.push(token === '(' ? group(true) : { token });
    }
    return terms;
  }
  function lookup(path) {
    if (path === '.') return context.dot;
    const parts = path.split('.');
    let value;
    if (path.startsWith('$.') || path === '$') { value = context.root; parts.shift(); }
    else if (path.startsWith('$')) { value = context.variables[parts.shift()]; }
    else if (path.startsWith('.')) { value = context.dot; parts.shift(); }
    else return path;
    for (const part of parts) {
      if (!part) continue;
      if (!value || ['__proto__', 'prototype', 'constructor'].includes(part)) return undefined;
      const parent = value;
      value = parent[part];
      if (typeof value === 'function') value = value.bind(parent);
    }
    return value;
  }
  const functions = {
    eq: (a, b) => a === b, ne: (a, b) => a !== b, lt: (a, b) => a < b, ge: (a, b) => a >= b,
    not: value => !truth(value), and: (...values) => values.every(truth), or: (...values) => values.some(truth),
    len: value => value?.length || 0, contains: (values, value) => values?.includes(value) || false,
  };
  function evaluate(node) {
    if (!Array.isArray(node)) {
      const token = node.token;
      if (token.startsWith('"')) return JSON.parse(token);
      if (/^-?\d+$/.test(token)) return Number(token);
      return lookup(token);
    }
    if (!node.length) return '';
    const name = node[0].token;
    const operation = functions[name];
    const head = operation || evaluate(node[0]);
    return typeof head === 'function' ? head(...node.slice(1).map(evaluate)) : head;
  }
  return evaluate(group());
}

const cache = new Map();
function parse(source) {
  if (cache.has(source)) return cache.get(source);
  const root = [];
  const stack = [];
  let current = root;
  for (const piece of source.split(/({{[\s\S]*?}})/g)) {
    if (!piece.startsWith('{{')) { current.push({ type: 'text', value: piece }); continue; }
    const directive = piece.slice(2, -2).trim().replace(/\s+/g, ' ');
    const [command] = directive.split(' ');
    if (['if', 'range', 'define'].includes(command)) {
      const node = { type: command, expression: directive.slice(command.length).trim(), children: [], branches: [] };
      if (command === 'if') node.branches.push({ expression: node.expression, children: node.children });
      current.push(node); stack.push({ node, parent: current }); current = node.children;
    } else if (command === 'else') {
      const frame = stack.at(-1);
      if (!frame) throw new Error('Unexpected template else');
      if (frame.node.type === 'if') {
        const branch = { expression: directive.startsWith('else if ') ? directive.slice(8) : null, children: [] };
        frame.node.branches.push(branch); current = branch.children;
      } else { frame.node.otherwise = []; current = frame.node.otherwise; }
    } else if (command === 'end') {
      const frame = stack.pop();
      if (!frame) throw new Error('Unexpected template end');
      current = frame.parent;
    } else if (command === 'template') current.push({ type: 'body' });
    else current.push({ type: 'value', expression: directive });
  }
  if (stack.length) throw new Error('Unclosed template directive');
  cache.set(source, root);
  return root;
}

export function renderTemplate(source, data) {
  function render(nodes, context) {
    return nodes.map(node => {
      if (node.type === 'text') return node.value;
      if (node.type === 'body') return context.root.Body?.value || '';
      if (node.type === 'define') return render(node.children, context);
      if (node.type === 'value') {
        const value = expression(node.expression, context);
        return value instanceof SafeHTML ? value.value : escapeHTML(value ?? '');
      }
      if (node.type === 'if') {
        const branch = node.branches.find(branch => branch.expression === null || truth(expression(branch.expression, context)));
        return branch ? render(branch.children, context) : '';
      }
      if (node.type === 'range') {
        const match = node.expression.match(/^(?:(\$\w+)(?:,\s*(\$\w+))?\s*:=\s*)?(.+)$/);
        const values = expression(match[3], context) || [];
        if (!values.length) return render(node.otherwise || [], context);
        return values.map((value, index) => {
          const variables = { ...context.variables };
          if (match[2]) { variables[match[1]] = index; variables[match[2]] = value; }
          else if (match[1]) variables[match[1]] = value;
          return render(node.children, { ...context, dot: value, variables });
        }).join('');
      }
      throw new Error(`Unsupported template node: ${node.type}`);
    }).join('');
  }
  return render(parse(source), { dot: data, root: data, variables: {} });
}
