import { categories } from '../data/categories.js';
import { reactionCount, commentCount } from '../services/feed.js';
import { plainPreview } from '../utils/content.js';

export const originalCategories = categories.map((category, index) => ({ ID: String(index + 1), Name: category.id }));

export function originalDate(iso) {
  return { Format(layout) {
    const date = new Date(iso);
    const day = date.toISOString().slice(0, 10);
    const time = date.toISOString().slice(11, 16);
    const month = date.toLocaleString('en', { month: 'short', timeZone: 'UTC' });
    if (layout === '2006-01-02 15:04') return `${day} ${time}`;
    if (layout === 'Jan 2, 15:04') return `${month} ${date.getUTCDate()}, ${time}`;
    return `${month} ${date.getUTCDate()}, ${date.getUTCFullYear()}`;
  } };
}

export function contentBlocks(content) {
  return content.split(/(```[\s\S]*?```)/g).filter(value => value.trim()).flatMap(part => {
    if (part.startsWith('```')) {
      const match = part.match(/^```([^\n]*)\n?([\s\S]*?)```$/);
      return [{ Type: 'code', Content: match[2].trimEnd(), Language: match[1].trim() }];
    }
    return part.split(/\n\s*\n/).filter(value => value.trim()).map(value => {
      const text = value.trim();
      if (/^https?:\/\/\S+$/i.test(text)) {
        try { const url = new URL(text); return { Type: 'link', Content: text, Link: { URL: url.href, Text: url.href } }; } catch { /* Render invalid URLs literally. */ }
      }
      return { Type: 'text', Content: text };
    });
  });
}

export function originalUser(user) {
  if (!user) return null;
  return { ID: user.id, Username: user.username, Email: `${user.username}@example.invalid`, Photo: user.photo || 'default.png' };
}

export function originalPost(post, state) {
  const author = state.users.find(user => user.id === post.authorId);
  const reaction = post.reactions[state.currentUserId];
  return {
    ID: post.id, Title: post.title, UserID: post.authorId, Username: author.username,
    Content: contentBlocks(post.content), Preview: plainPreview(post.content),
    CreatedAt: originalDate(post.createdAt), Category: post.categories.join(', '),
    LikeCount: reactionCount(post, 'like'), DislikeCount: reactionCount(post, 'dislike'),
    UserLiked: reaction === 'like', UserDisliked: reaction === 'dislike', CommentCount: commentCount(state, post.id),
    Comments: state.comments.filter(comment => comment.postId === post.id).map(comment => ({
      ID: comment.id, PostID: post.id, Username: state.users.find(user => user.id === comment.authorId).username,
      Content: comment.content, CreatedAt: originalDate(comment.createdAt), LikeCount: reactionCount(comment, 'like'),
      DislikeCount: reactionCount(comment, 'dislike'), UserLiked: comment.reactions[state.currentUserId] === 'like',
      UserDisliked: comment.reactions[state.currentUserId] === 'dislike',
    })),
  };
}

export function originalTheme(category) {
  const theme = {
    Accent: category.color, AccentDark: category.dark, AccentLight: category.light, Secondary: category.secondary,
    BgPrimary: '#0a0e1a', BgSecondary: '#1a1f2e', BgCard: '#242938', BgElevated: '#2d3748',
    TextPrimary: '#f7fafc', TextSecondary: '#cbd5e0', TextMuted: '#a0aec0', Border: '#4a5568', BorderLight: '#2d3748',
    Shadow: '0 10px 25px rgba(0, 212, 255, 0.1)', ShadowStrong: '0 20px 40px rgba(0, 212, 255, 0.2)',
    BoxShadow: '0 4px 12px rgba(255,255,255,0)',
    Radius: '16px', RadiusSmall: '8px', Spacing: '24px',
  };
  // Exact per-language values from Internal/api/categoryhand.go, including
  // the original non-accent shadow hues and Brainfuck's dark highlight.
  const shadows = {
    python: ['55,118,171', '.18', '.30'], javascript: ['247,223,30', '.15', '.28'],
    rust: ['206,66,43', '.15', '.28'], java: ['255,79,79', '.15', '.30'],
    ruby: ['215,0,55', '.15', '.30'], matlab: ['255,119,14', '.15', '.30'], brainfuck: ['255,106,213', '.15', '.30'],
  }[category.id];
  if (shadows) {
    theme.Shadow = `0 10px 25px rgba(${shadows[0]},${shadows[1]})`;
    theme.ShadowStrong = `0 20px 40px rgba(${shadows[0]},${shadows[2]})`;
    theme.BoxShadow = `0 4px 12px rgba(${category.id === 'javascript' ? '255,233,0' : shadows[0]},.3)`;
  }
  if (category.id === 'brainfuck') theme.AccentLight = '#0d0b14';
  return theme;
}
