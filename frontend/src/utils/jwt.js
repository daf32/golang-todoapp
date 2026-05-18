export function parseJwt(token) {
  if (!token) return null;
  try {
    const b64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(atob(b64));
  } catch {
    return null;
  }
}

export function getUserIdFromToken(token) {
  const claims = parseJwt(token);
  return claims?.sub != null ? Math.round(Number(claims.sub)) : null;
}

export function getRoleFromToken(token) {
  return parseJwt(token)?.role ?? null;
}
