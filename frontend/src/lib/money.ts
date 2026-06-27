/**
 * Format an integer number of cents as a USD string.
 * centsToUSD(1999) => "$19.99"
 */
export function centsToUSD(cents: number): string {
  const value = (cents ?? 0) / 100;
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(value);
}
