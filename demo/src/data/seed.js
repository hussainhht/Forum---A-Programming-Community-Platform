import discussions from './discussions.json' with { type: 'json' };
import { categories, avatarColors } from './categories.js';

const profiles = [
  ['demo_dev', 'Alex Byte', 'Exploring the forum, one good question at a time.', ['golang', 'javascript']],
  ['channel_chaser', 'Mira Loop', 'Building small Go tools and learning how to make concurrency boring.', ['golang', 'rust']],
  ['async_owl', 'Ellis Finch', 'Browser tinkerer. I like accessible interfaces and understandable promises.', ['javascript', 'python']],
  ['borrow_checker', 'Rowan Moss', 'Learning systems programming with a notebook full of ownership diagrams.', ['rust', 'golang']],
  ['ruby_spark', 'Cleo Reed', 'Writing tiny Ruby scripts that make everyday work a little easier.', ['ruby', 'python']],
  ['list_comprehender', 'Remy Cloud', 'Data wrangler, generator enthusiast, and friendly code reviewer.', ['python', 'matlab']],
  ['matrix_maven', 'Nova Vale', 'Turning noisy signals into plots. Always checking the matrix dimensions.', ['matlab', 'python']],
  ['tape_walker', 'Kit Pixel', 'Esolang puzzle collector. My favorite debugger is a carefully drawn tape.', ['brainfuck', 'rust']],
  ['java_journeys', 'Sage Maple', 'Learning object design by building and rebuilding small Java projects.', ['java', 'golang']],
  ['stack_sprout', 'Ari Cedar', 'New to programming. Asking the questions I used to be afraid to ask.', ['javascript', 'java']],
  ['mutex_moon', 'Jules Orbit', 'Interested in race detectors, practical testing, and clear error messages.', ['golang', 'rust']],
  ['query_quill', 'Robin Slate', 'Database experiments, query plans, and a healthy respect for backups.', ['golang', 'python']],
  ['lambda_lark', 'Indigo Wren', 'I enjoy finding the simplest explanation for a complicated loop.', ['ruby', 'javascript']],
  ['pixel_parser', 'Tess Ember', 'Building parsers and little tools for text. Fictional bugs, real curiosity.', ['rust', 'brainfuck']],
  ['heap_hiker', 'Finley Brook', 'Measuring allocations before making performance claims.', ['java', 'rust']],
  ['vector_voyager', 'Luca Comet', 'Numerical computing student interested in stable algorithms.', ['matlab', 'python']],
  ['byte_blossom', 'Wren Echo', 'Weekend Ruby projects and readable code over clever code.', ['ruby', 'java']],
  ['null_navigator', 'Piper Arc', 'Frontend explorer collecting useful patterns and smaller bug reports.', ['javascript', 'brainfuck']],
];

const followUps = [
  'The producer-owned close was the missing detail. I traced the range loop and can now explain why it stops after the last value.',
  'Returning after the JSON error was exactly what I missed in my handler. I have added a malformed-input case to the test.',
  'I will compare the join plans with representative category sizes before adding another index. The insert cost belongs in that comparison too.',
  'Delegating from a stable feed element fixed the newly inserted cards. Clicking the nested icon now reaches the same handler as clicking the label.',
  'I will keep ISO date strings at the storage boundary and use structuredClone only for in-memory values that actually need those richer types.',
  'The borrowed &str signature makes the ownership boundary much clearer. I no longer need to clone the message just to print it.',
  'I had treated reserved slots as initialized values. Separating length from capacity explains both the panic and why push is still needed.',
  'I will compare release builds and make the result observable. The loop version is staying until I have evidence to justify changing it.',
  'Assigning the new array inside the hash default block was the part I overlooked. A second access to the same key now keeps its previous filenames.',
  'Keeping JSON keys as strings is the simplest choice for this importer. I have added explicit checks for the required theme setting instead of relying on key conversion.',
  'A single comprehension is readable here. For the reporting code with several branches, I am keeping a normal loop and named helper.',
  'Keeping the CSV reader lazy made the memory use predictable. I have also started reporting the row number when a quantity cannot be parsed.',
  'I was reusing the generator object instead of calling the function again. The two-pass example now creates independent iterators.',
  'The dimensions are now fixed before the simulation starts. I will use single precision only after checking whether its accuracy is sufficient.',
  'Applying the same logical mask to timestamps prevented a subtle alignment bug in my chart. I am also rejecting infinite readings explicitly.',
  'Tracing the pointer after each loop iteration made the final >+ obvious. I added the unsigned-byte and output assumptions beside the example.',
  'The echo loop depends on the interpreter mapping EOF to zero. I will put that assumption in the example rather than calling it portable.',
  'The destination must start at zero for this to be a move. With a nonzero destination it is an addition, which explains my unexpected tape state.',
  'The record example gives me a compact immutable key. I will keep mutable display details outside the object used for HashSet membership.',
  'ArrayDeque matches the actual queue operations much better. I was comparing insertion costs while ignoring the traversal needed to reach the index.',
];

// Stable fictional identities, dates, references and per-user reactions.
// Each post has two subject-specific replies; every other post also has an
// author follow-up. No network calls or randomness occur during seeding.
export function createSeed() {
  const users = profiles.map(([username, displayName, bio, interests], index) => ({
    id: `user-${index + 1}`, username, displayName, bio, interests,
    avatar: avatarColors[index % avatarColors.length],
    createdAt: new Date(Date.UTC(2025, index % 12, index + 1)).toISOString(),
  }));
  const reactions = (index, count) => Object.fromEntries(Array.from({ length: count }, (_, j) => [
    users[1 + (index + j) % (users.length - 1)].id, j === 0 && index % 4 === 0 ? 'dislike' : 'like',
  ]));
  const posts = discussions.map(([category, title, body, code], index) => ({
    id: String(index + 1), authorId: users[1 + index % (users.length - 1)].id,
    title, content: `${body}\n\n\`\`\`${index === 4 ? 'sql' : categories.find(c => c.id === category).language}\n${code}\n\`\`\``,
    categories: index === 4 ? [category, 'python'] : [category],
    createdAt: new Date(Date.UTC(2026, 9, 26 - Math.floor((index % 5 * 8 + Math.floor(index / 5)) / 2), 18 - index % 12, index % 60)).toISOString(),
    reactions: reactions(index, 2 + index * 7 % 15),
  }));
  const comments = discussions.flatMap((discussion, index) => {
    const answers = [discussion[4], discussion[5]];
    if (index % 2 === 0) answers.push(followUps[index / 2]);
    return answers.map((content, reply) => ({
      id: `comment-${index + 1}-${reply + 1}`, postId: posts[index].id,
      authorId: reply === 2 ? posts[index].authorId : users[1 + (index + reply + 3) % (users.length - 1)].id,
      content: index === 0 && reply === 0 ? `${content}\n\n\`\`\`go\nfor value := range values {\n    fmt.Println(value)\n}\n\`\`\`` : content,
      createdAt: new Date(Date.parse(posts[index].createdAt) + (reply + 1) * 3600000).toISOString(),
      reactions: reactions(index + reply, (index + reply) % 7),
    }));
  });
  return { version: 1, users, posts, comments, currentUserId: null };
}
