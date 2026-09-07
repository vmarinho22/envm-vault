/**
 * Turns a display name into a URL-safe slug: lowercase ASCII,
 * accents stripped, non-alphanumeric runs collapsed to a single hyphen.
 * Must stay in sync with api/internal/slug.FromName.
 */
export function fromName(name: string): string {
  const trimmed = name.trim();
  if (trimmed === '') {
    return '';
  }

  return trimmed
    .normalize('NFD')
    .replace(/\p{M}/gu, '')
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}
