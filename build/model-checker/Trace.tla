------------------------- MODULE Trace -------------------------

INSTANCE JSONUtils

Trace == JSONDeserialize("/etc/model-checker/data/trace.json")

=============================================================================
