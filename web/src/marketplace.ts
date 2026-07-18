import type { CatalogGroup, CatalogModel } from '~/src/api'

// ----------------------------------------------------------------------------
// Model square (模型广场) shared logic — ported from xinghai-api pricing feature
// ----------------------------------------------------------------------------

export type TokenUnit = 'M' | 'K'
export type ViewMode = 'card' | 'table'
export type SortOption = 'name' | 'price-low' | 'price-high'
export type PriceKind = 'input' | 'output' | 'cache'

export const FILTER_ALL = 'all'
export const PAGE_SIZE = 20

/** A catalog entry augmented with derived vendor metadata. */
export interface SquareModel extends CatalogModel {
  vendor_name: string
  vendor_slug: string
  /** Reserved for future backend-provided catalog descriptions. */
  description?: string
}

export function toSquareModel(item: CatalogModel): SquareModel {
  return { ...item, vendor_name: item.provider, vendor_slug: item.provider_slug }
}

/** Unique vendors across the catalog, with model counts, preserving order. */
export function extractVendors(models: SquareModel[]) {
  const map = new Map<string, { name: string; slug: string; count: number }>()
  for (const model of models) {
    const entry = map.get(model.vendor_name)
    if (entry) entry.count += 1
    else map.set(model.vendor_name, { name: model.vendor_name, slug: model.vendor_slug, count: 1 })
  }
  return [...map.values()]
}

/** Lowest-multiplier group of a model, or the selected group when filtering. */
export function getDisplayGroup(model: CatalogModel, selectedGroup: string): CatalogGroup | undefined {
  const groups = model.groups
  if (!groups.length) return undefined
  if (selectedGroup !== FILTER_ALL) {
    const hit = groups.find((group) => group.id === selectedGroup)
    if (hit) return hit
  }
  let best = groups[0]
  for (const group of groups.slice(1)) {
    if (Number(group.multiplier) < Number(best.multiplier)) best = group
  }
  return best
}

function basePrice(model: CatalogModel, kind: PriceKind): number | null {
  const value = kind === 'input' ? model.input_per_million : kind === 'output' ? model.output_per_million : model.cached_input_per_million
  if (value == null) return null
  const numeric = Number(value)
  // The backend stores an omitted cached-input price as 0 — treat it as unconfigured.
  if (kind === 'cache' && numeric === 0) return null
  return numeric
}

/** Effective local-currency price per 1M tokens for the display group. */
export function effectivePrice(model: CatalogModel, kind: PriceKind, selectedGroup: string): number | null {
  const base = basePrice(model, kind)
  if (base == null) return null
  const modelMultiplier = Number(model.multiplier ?? 1)
  const group = getDisplayGroup(model, selectedGroup)
  return base * modelMultiplier * Number(group?.multiplier ?? 1)
}

/** Effective price for a concrete group (used by the group pricing table). */
export function groupPrice(model: CatalogModel, kind: PriceKind, group: CatalogGroup): number | null {
  const base = basePrice(model, kind)
  if (base == null) return null
  return base * Number(model.multiplier ?? 1) * Number(group.multiplier ?? 1)
}

/** Format a per-1M price for the chosen token unit, trimming trailing zeros. */
export function formatSquarePrice(value: number | null, unit: TokenUnit): string {
  if (value == null) return ''
  const scaled = unit === 'K' ? value / 1000 : value
  if (scaled === 0) return '¥0'
  const digits = scaled >= 100 ? 2 : scaled >= 1 ? 4 : 6
  const fixed = scaled.toFixed(digits)
  const trimmed = String(Number.parseFloat(fixed))
  return `¥${trimmed}`
}

/** Format a group multiplier like x1 / x0.5 / x1.25. */
export function formatRatio(multiplier: number): string {
  const value = Number(multiplier)
  const formatted = Number.isInteger(value) ? String(value) : value.toFixed(3).replace(/0+$/, '').replace(/\.$/, '')
  return `x${formatted}`
}

/** Deterministic accent color for vendor badges, like StatusBadge autoColor. */
export function vendorColor(name: string): { bg: string; fg: string } {
  let hash = 0
  for (let index = 0; index < name.length; index += 1) hash = (hash * 31 + name.charCodeAt(index)) >>> 0
  const hue = hash % 360
  return { bg: `hsl(${hue} 65% 45% / 0.14)`, fg: `hsl(${hue} 55% 55%)` }
}

/** Lobehub static icon for a provider slug (same source as provider admin view). */
export function vendorIconUrl(slug: string): string {
  return `https://unpkg.com/@lobehub/icons-static-svg@1.93.0/icons/${slug}.svg`
}

export interface SquareFilters {
  search: string
  vendor: string
  group: string
  quotaType: string
  sortBy: SortOption
}

export function filterAndSort(models: SquareModel[], filters: SquareFilters): SquareModel[] {
  const query = filters.search.trim().toLowerCase()
  let result = models
  if (query) {
    result = result.filter((model) =>
      model.model.toLowerCase().includes(query) || model.vendor_name.toLowerCase().includes(query),
    )
  }
  if (filters.vendor !== FILTER_ALL) result = result.filter((model) => model.vendor_name === filters.vendor)
  if (filters.group !== FILTER_ALL) result = result.filter((model) => model.groups.some((group) => group.id === filters.group))
  const sorted = [...result]
  if (filters.sortBy === 'name') sorted.sort((a, b) => a.model.localeCompare(b.model))
  else if (filters.sortBy === 'price-low') sorted.sort((a, b) => (effectivePrice(a, 'input', filters.group) ?? Number.POSITIVE_INFINITY) - (effectivePrice(b, 'input', filters.group) ?? Number.POSITIVE_INFINITY))
  else if (filters.sortBy === 'price-high') sorted.sort((a, b) => (effectivePrice(b, 'input', filters.group) ?? -1) - (effectivePrice(a, 'input', filters.group) ?? -1))
  return sorted
}
