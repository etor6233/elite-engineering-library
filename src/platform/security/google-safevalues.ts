/**
 * Google SafeValues is the admitted DOM-XSS primitive for non-React DOM sinks.
 * React text interpolation remains preferred; raw HTML still requires an explicit
 * product policy, provenance and review before calling a sanitizer.
 */
export {
  htmlEscape,
  isHtml,
  sanitizeHtml,
  sanitizeHtmlAssertUnchanged,
  unwrapHtml,
} from "safevalues";

export {
  setAnchorHref,
  setElementInnerHtml,
  setIframeSrcdoc,
  setLocationHref,
  setScriptSrc,
} from "safevalues/dom";
