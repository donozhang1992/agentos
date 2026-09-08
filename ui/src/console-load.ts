// Copyright 2026 Dapeng Zhang and Agenova contributors.
// SPDX-License-Identifier: Apache-2.0
import type { ClaimRequest, IssuedState } from './contracts.generated';
import type { EvidenceResult, EvidenceSource } from './evidence-source';
import { issuedKey, requestKey, type ConsoleRoute } from './console-route';
import { shapeDiagnostics } from './shape-check';

// Local composition only; never serialized or offered as a live evidence DTO.
export type ConsoleResult =
  | { status: 'ready'; issued: IssuedState; request?: ClaimRequest; requestGap?: string }
  | Exclude<EvidenceResult, { status: 'request' | 'issued' }>;

export async function loadConsole(source: EvidenceSource, route: ConsoleRoute): Promise<ConsoleResult> {
  try {
    const result = await source.load(issuedKey(route));
    if (result.status === 'request') return { status: 'invalid', diagnostics: [{ category: 'unexpected-source-kind', fieldPath: '$' }] };
    if (result.status !== 'issued') return result;
    const state = result.data;
    const diagnostics = shapeDiagnostics('IssuedState', state);
    if (diagnostics.length) return { status: 'invalid', diagnostics };
    if ((route.kind === 'requests' ? state.requestRef !== route.reference : state.claim?.id !== route.reference) ||
        state.evidence.requestRef !== state.requestRef || (state.claim &&
        (state.claim.requestRef !== state.requestRef || state.claim.templateRef !== state.action.templateRef))) {
      return { status: 'invalid', diagnostics: [{ category: 'reference-mismatch', fieldPath: 'requestRef' }] };
    }
    let requestResult: EvidenceResult;
    try { requestResult = await source.load(requestKey(state.requestRef)); }
    catch { return { status: 'ready', issued: state, requestGap: 'Matching request document unavailable.' }; }
    if (requestResult.status !== 'request') return {
      status: 'ready', issued: state,
      requestGap: requestResult.status === 'not-found' ? 'Matching request document not supplied.' : 'Matching request document unavailable or malformed.',
    };
    const request = requestResult.data;
    if (shapeDiagnostics('ClaimRequest', request).length || request.metadata.name !== state.requestRef || request.spec.templateRef !== state.action.templateRef) {
      return { status: 'ready', issued: state, requestGap: 'Request correlation failed. Access comparison withheld.' };
    }
    return { status: 'ready', issued: state, request };
  } catch { return { status: 'unavailable' }; }
}
