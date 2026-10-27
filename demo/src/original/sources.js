// Import the real application, rather than maintaining a second design.
const templates = import.meta.glob('../../../web/templates/*.html', { eager: true, query: '?raw', import: 'default' });
const styles = import.meta.glob('../../../web/static/css/*.css', { eager: true, query: '?raw', import: 'default' });
export const template = name => templates[`../../../web/templates/${name}.html`];

export function originalCSS(names) {
  return names.map(name => {
    const source = styles[`../../../web/static/css/${name}.css`];
    if (!source) throw new Error(`Missing original stylesheet: ${name}`);
    // The original Inter import is redundant: root.css supplies the forum's
    // system font. Keep the demo's application assets entirely local.
    return source.replace(/@import\s+url\([^)]+fonts\.googleapis[^)]+\);?/g, '')
      .replace(/url\((['"]?)\.\.\//g, `url($1${import.meta.env.BASE_URL}web/static/`)
      .replace(/bg-tech-purple\.png/g, 'bg-tech.png');
  }).join('\n');
}
