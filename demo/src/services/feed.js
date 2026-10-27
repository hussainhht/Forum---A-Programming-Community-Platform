export const reactionCount = (item, type) => Object.values(item.reactions).filter(value => value === type).length;
export const commentCount = (state, postId) => state.comments.filter(comment => comment.postId === postId).length;

export function selectPosts(state, { query = '', categories = [], sort = 'latest', authorId } = {}) {
  const search = query.trim().toLocaleLowerCase();
  const posts = state.posts.filter(post => (!authorId || post.authorId === authorId)
    && (!categories.length || categories.some(id => post.categories.includes(id)))
    && (!search || `${post.title}\n${post.content}`.toLocaleLowerCase().includes(search)));
  const latest = (a, b) => Date.parse(b.createdAt) - Date.parse(a.createdAt) || a.id.localeCompare(b.id);
  return posts.sort((a, b) => {
    if (sort === 'popular') return reactionCount(b, 'like') - reactionCount(a, 'like') || latest(a, b);
    if (sort === 'comments') return commentCount(state, b.id) - commentCount(state, a.id) || latest(a, b);
    return latest(a, b);
  });
}

export function paginate(items, page = 1, perPage = 9) {
  const pages = Math.max(1, Math.ceil(items.length / perPage));
  const current = Math.min(pages, Math.max(1, Number.isFinite(Number(page)) ? Math.floor(Number(page)) : 1));
  return { items: items.slice((current - 1) * perPage, current * perPage), page: current, pages, total: items.length };
}
