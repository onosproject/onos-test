------------------------- MODULE Trace -------------------------

INSTANCE IOUtils

Trace == IODeserialize("/etc/model-checker/data/trace.bin", FALSE)

=============================================================================
