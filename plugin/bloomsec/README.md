# bloomsec

## Name

*bloomsec* - adds DNSSEC records and Bloom filter TXT records to zone files.

## Description

The *bloomsec* plugin is used to add a Bloom filter encoding the zone's data and to sign (see RFC 6781) zones. In this process TXT resource records containing chunks of the Bloom filter and DNSSEC resource records are
added. The signatures that sign the resource records sets have an expiration date, this means the
signing process must be repeated before this expiration data is reached. Otherwise the zone's data
will go BAD (RFC 4035, Section 5.5). The *sign* plugin takes care of this.

Only NSEC is supported, *bloomsec* does *not* support NSEC3.

*bloomsec* works in conjunction with the *bloomfile* plugins; this plugin **creates** the Bloom filter and **signs** the zones
files, and *bloomfile* **serve** the zones *data*.

For this plugin to work at least one Common Signing Key, (see coredns-keygen(1)) is needed. This key
(or keys) will be used to sign the entire zone. *bloomsec* does *not* support the ZSK/KSK split, nor will
it do key or algorithm rollovers - it just signs.

*bloomsec* will:
 * Create a Bloom filter encoding all the domain names for which the zone is authoritative. Furthermore, it chunks it into equally sized chunks and puts the into TXT resource records.

 *  (Re)-sign the zone with the CSK(s) when:

     -  the last time it was signed is more than a 6 days ago. Each zone will have some jitter
        applied to the inception date.

     -  the signature only has 14 days left before expiring.

    Both these dates are only checked on the SOA's signature(s).

 *  Create RRSIGs that have an inception of -3 hours (minus a jitter between 0 and 18 hours)
    and a expiration of +32 (plus a jitter between 0 and 5 days) days for every given DNSKEY.

 *  Add NSEC records for all names in the zone. The TTL for these is the negative cache TTL from the
    SOA record.

 *  Add or replace *all* apex CDS/CDNSKEY records with the ones derived from the given keys. For
    each key two CDS are created one with SHA1 and another with SHA256.

 *  Update the SOA's serial number to the *Unix epoch* of when the signing happens. This will
    overwrite *any* previous serial number.


There are two ways that dictate when a zone is signed. Normally every 6 days (plus jitter) it will
be resigned. If for some reason we fail this check, the 14 days before expiring kicks in.

Keys are named (following BIND9): `K<name>+<alg>+<id>.key` and `K<name>+<alg>+<id>.private`.
The keys **must not** be included in your zone; they will be added by *bloomsec*. These keys can be
generated with `coredns-keygen` or BIND9's `dnssec-keygen`. You don't have to adhere to this naming
scheme, but then you need to name your keys explicitly, see the `keys file` directive.

A generated zone is written out in a file named `db.<name>.signed` in the directory named by the
`directory` directive (which defaults to `/var/lib/coredns`).

## Syntax

~~~
bloomsec DBFILE [ZONES...] {
    key file|directory KEY...|DIR...
    directory DIR
    fp_rate FPRATE
    chunksize CHUNKSIZE
}
~~~

*  **DBFILE** the zone database file to read and parse. If the path is relative, the path from the
   *root* plugin will be prepended to it.
*  **ZONES** zones it should be sign for. If empty, the zones from the configuration block are
   used.
* `key` specifies the key(s) (there can be multiple) to sign the zone. If `file` is
   used the **KEY**'s filenames are used as is. If `directory` is used, *sign* will look in **DIR**
   for `K<name>+<alg>+<id>` files. Any metadata in these files (Activate, Publish, etc.) is
   *ignored*. These keys must also be Key Signing Keys (KSK).
*  `directory` specifies the **DIR** where CoreDNS should save zones that have been signed.
   If not given this defaults to `/var/lib/coredns`. The zones are saved under the name
   `db.<name>.signed`. If the path is relative the path from the *root* plugin will be prepended
   to it.
* `fp_rate` specifies the **FPRATE** which should be guaranteed by the Bloom filter. This parameter must be a floating point number between (excluding) 0 and 1.
    If not given this defaults to `0.001`.
* `chunksize` specifies the **CHUNKSIZE** in bits. This is the number of bits included in the data section of the TXT resource record. This parameter is an integer larger or equal than 2040 that needs to be a multiple of 2040.
    If not given this defaults to `28560`.

Keys can be generated with `coredns-keygen`, to create one for use in the *sign* plugin, use:
`coredns-keygen example.org` or `dnssec-keygen -a ECDSAP256SHA256 -f KSK example.org`.

## Examples

Create the Bloom filter and sign the `example.org` zone contained in the file `db.example.org` and write the result to
`./db.example.org.signed` to let the *bloomfile* plugin pick it up and serve it. The keys used
are read from `/etc/coredns/keys/Kexample.org.key` and `/etc/coredns/keys/Kexample.org.private`. The chunk size and the false positive rate use their default values.

~~~ txt
example.org {
    bloomfile db.example.org.signed

    bloomsec db.example.org {
        key file /etc/coredns/keys/Kexample.org
        directory .
    }
}
~~~

Running this leads to the following log output (note the timers in this example have been set to
shorter intervals).

~~~ txt
[WARNING] plugin/bloomfile: Failed to open "open /tmp/db.example.org.signed: no such file or directory": trying again in 1m0s
[INFO] plugin/bloomsec: Signing "example.org." because open /tmp/db.example.org.signed: no such file or directory
[INFO] plugin/bloomsec: Successfully signed zone "example.org." in "/tmp/db.example.org.signed" with key tags "59725" and 1564766865 SOA serial, elapsed 9.357933ms, next: 2019-08-02T22:27:45.270Z
[INFO] plugin/bloomfile: Successfully reloaded zone "example.org." in "/tmp/db.example.org.signed" with serial 1564766865
~~~

Or use a single zone file for *multiple* zones, note that the **ZONES** are repeated for both plugins.
Also note this outputs *multiple* signed output files. Here we use the default output directory
`/var/lib/coredns`.

~~~ txt
. {
    bloomfile /var/lib/coredns/db.example.org.signed example.org
    bloomfile /var/lib/coredns/db.example.net.signed example.net
    bloomsec db.example.org example.org example.net {
        key directory /etc/coredns/keys
    }
}
~~~

This is the same configuration, but the zones are put in the server block, but note that you still
need to specify what file is served for what zone in the *file* plugin:

~~~ txt
example.org example.net {
    bloomfile var/lib/coredns/db.example.org.signed example.org
    bloomfile var/lib/coredns/db.example.net.signed example.net
    bloomsec db.example.org {
        key directory /etc/coredns/keys
    }
}
~~~

Be careful to fully list the origins you want to sign, if you don't:

~~~ txt
example.org example.net {
    bloomsec plugin/sign/testdata/db.example.org miek.org {
        key file /etc/coredns/keys/Kexample.org
    }
}
~~~

This will lead to `db.example.org` be signed *twice*, as this entire section is parsed twice because
you have specified the origins `example.org` and `example.net` in the server block.

Forcibly resigning a zone can be accomplished by removing the signed zone file (CoreDNS will keep
on serving it from memory), and sending SIGUSR1 to the process to make it reload and resign the zone
file.

## See Also

The DNSSEC RFCs: RFC 4033, RFC 4034 and RFC 4035. And the BCP on DNSSEC, RFC 6781. Further more the
manual pages coredns-keygen(1) and dnssec-keygen(8). And the *file* plugin's documentation.

Coredns-keygen can be found at
[https://github.com/coredns/coredns-utils](https://github.com/coredns/coredns-utils) in the
coredns-keygen directory.

Other useful DNSSEC tools can be found in [ldns](https://nlnetlabs.nl/projects/ldns/about/), e.g.
`ldns-key2ds` to create DS records from DNSKEYs.

## Bugs

`keys directory` is not implemented.
