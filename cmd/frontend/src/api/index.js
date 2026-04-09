const BASE = ''

async function request(path, options = {}) {
  const res = await fetch(BASE + path, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || 'Something went wrong')
  return data
}

export const api = {
  signUp: (body) => request('/auth/sign-up', { method: 'POST', body: JSON.stringify(body) }),
  signIn: (body) => request('/auth/sign-in', { method: 'POST', body: JSON.stringify(body) }),

  getEntries: () => request('/entries/'),
  getEntry: (id) => request(`/entries/${id}`),
  createEntry: (body) => request('/entries/', { method: 'POST', body: JSON.stringify(body) }),
  deleteEntry: (id) => request(`/entries/${id}`, { method: 'DELETE' }),
  getInsight: () => request('/entries/insight'),
}
