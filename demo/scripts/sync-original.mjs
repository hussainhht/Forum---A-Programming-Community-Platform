import { cp, mkdir, copyFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const demo = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const repository = path.dirname(demo);
await mkdir(path.join(demo, 'public', 'web'), { recursive: true });
await cp(path.join(repository, 'web', 'static'), path.join(demo, 'public', 'web', 'static'), { recursive: true });
await mkdir(path.join(demo, 'public', 'profile-img'), { recursive: true });
await copyFile(path.join(repository, 'web', 'static', 'img', 'default.png'), path.join(demo, 'public', 'profile-img', 'default.png'));
// Fill assets referenced by the original templates but absent in the source.
await copyFile(path.join(repository, 'web', 'static', 'icons', 'info.svg'), path.join(demo, 'public', 'web', 'static', 'icons', 'Terms.svg'));
const { writeFile } = await import('node:fs/promises');
await writeFile(path.join(demo, 'public', 'web', 'static', 'icons', 'categories', 'matlab.svg'), '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32"><path fill="#ff9f40" d="m2 23 9-5 7-16 4 17 8 9-12-4-7 5z"/></svg>');
await writeFile(path.join(demo, 'public', 'web', 'static', 'icons', 'categories', 'brainfuck.svg'), '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32"><path fill="none" stroke="#ff6ad5" stroke-width="3" d="M10 5H5v22h5M22 5h5v22h-5M12 16h8M16 12v8"/></svg>');
console.log('Original forum assets prepared. Source files are unchanged.');
