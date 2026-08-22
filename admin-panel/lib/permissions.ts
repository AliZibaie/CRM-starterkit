export function hasPermission(userPerms: string[], key: string): boolean {
  return userPerms.includes(key);
}

export function hasAnyPermission(userPerms: string[], keys: string[]): boolean {
  if (keys.length === 0) return true;
  return keys.some((key) => userPerms.includes(key));
}
