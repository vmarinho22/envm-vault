export interface Tag {
  id: string;
  title: string;
  color: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTagInput {
  title: string;
  color: string;
}

export interface UpdateTagInput {
  title?: string;
  color?: string;
}

export interface ListTagsParams {
  title?: string;
}

export const DEFAULT_TAG_COLOR = '#6366F1';

/** Random saturated mid-tone hex for new tags (`#RRGGBB`). */
export function randomTagColor(): string {
  const hue = Math.floor(Math.random() * 360);
  const saturation = 0.55 + Math.random() * 0.3;
  const lightness = 0.4 + Math.random() * 0.2;
  return hslToHex(hue, saturation, lightness);
}

function hslToHex(h: number, s: number, l: number): string {
  const c = (1 - Math.abs(2 * l - 1)) * s;
  const x = c * (1 - Math.abs(((h / 60) % 2) - 1));
  const m = l - c / 2;
  let r = 0;
  let g = 0;
  let b = 0;
  if (h < 60) {
    r = c;
    g = x;
  } else if (h < 120) {
    r = x;
    g = c;
  } else if (h < 180) {
    g = c;
    b = x;
  } else if (h < 240) {
    g = x;
    b = c;
  } else if (h < 300) {
    r = x;
    b = c;
  } else {
    r = c;
    b = x;
  }
  const toHex = (n: number) =>
    Math.round((n + m) * 255)
      .toString(16)
      .padStart(2, '0')
      .toUpperCase();
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`;
}
