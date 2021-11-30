load("@ytt:struct", "struct")
load("@ytt:assert", "assert")
load("@ytt:yaml", "yaml")
load("_profiles.star", "profiles")

def test_merge_ingress_values():
    cases = [
        {
            'pkg_values': struct.encode({'someKey': 'someValue'}),
            'ingress_values': struct.encode({
                'ingressEnabled': True,
                'ingressDomain': 'example.com',
                'tls': {
                    'namespace': 'my-namespace',
                    'secretName': 'my-secret',
                },
            }),
            'expected': lambda result: result.someKey == 'someValue',
            'message': 'expected pkg_values to be preserved',
        },
        {
            'pkg_values': struct.encode({'someKey': 'someValue'}),
            'ingress_values': struct.encode({
                'ingressEnabled': True,
                'ingressDomain': 'example.com',
                'tls': {
                    'namespace': 'my-namespace',
                    'secretName': 'my-secret',
                },
            }),
            'expected': lambda result: (result.ingressEnabled == True and
                                        result.ingressDomain == 'example.com' and
                                        result.tls.namespace == 'my-namespace' and
                                        result.tls.secretName == 'my-secret'),
            'message': 'expected ingress_values to be passed through to the pkg_values',
        },
    ]

    for case in cases:
        result = profiles.merge_ingress_values(case['pkg_values'], case['ingress_values'])

        if not case['expected'](result):
            print(struct.decode(result))
            assert.fail(case['message'])
        end
    end
end

test_merge_ingress_values()
