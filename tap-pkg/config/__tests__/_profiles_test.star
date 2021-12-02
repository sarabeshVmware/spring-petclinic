load("@ytt:struct", "struct")
load("@ytt:assert", "assert")
load("@ytt:yaml", "yaml")
load("_profiles.star", "profiles")

def test_merge_ingress_values():
    cases = [
        {
            'short_description': 'preserve other pkg_values',
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
            'failure_message': 'expected pkg_values to be preserved',
        },
        {
            'short_description': 'pass through ingress vals',
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
            'failure_message': 'expected ingress_values to be passed through to the pkg_values',
        },
        {
            'short_description': 'pkg precedence',
            'pkg_values': struct.encode({
                'someKey': 'someValue',
                'ingressEnabled': False,
                'ingressDomain': 'nothing.example.com',
                'tls': {
                    'namespace': 'some-other-namespace',
                    'secretName': 'some-other-secret',
                },
            }),
            'ingress_values': struct.encode({
                'ingressEnabled': True,
                'ingressDomain': 'example.com',
                'tls': {
                    'namespace': 'my-namespace',
                    'secretName': 'my-secret',
                },
            }),
            'expected': lambda result: (result.ingressEnabled == False and
                                        result.ingressDomain ==  'nothing.example.com' and
                                        result.tls.namespace == 'some-other-namespace' and
                                        result.tls.secretName == 'some-other-secret'),
            'failure_message': 'expected pkg_values to take precedence over ingress_values',
        },
        {
            'short_description': 'ingress values omitted',
            'pkg_values': struct.encode({
                'someKey': 'someValue',
            }),
            'ingress_values': struct.encode({}),
            'expected': lambda result: (lambda result_dict:
                                        'ingressEnabled' not in result_dict and
                                        'ingressDomain' not in result_dict and
                                        'tls' not in result_dict)(struct.decode(result)),
            'failure_message': 'omitted ingress values should not change pkg_values',
        },
        {
            'short_description': 'empty ingress tls dict',
            'pkg_values': struct.encode({
                'someKey': 'someValue',
                'tls': {
                    'namespace': 'some-other-namespace',
                    'secretName': 'some-other-secretName',
                },
            }),
            'ingress_values': struct.encode({
                'tls': {},
            }),
            'expected': lambda result: (result.tls.namespace == 'some-other-namespace' and
                                        result.tls.secretName == 'some-other-secretName'),
            'failure_message': 'expected empty ingress tls dict to have no effect on pkg_values',
        },
    ]

    for case in cases:
        if 'short_description' in case:
            print('Running:', case['short_description'])
        end
        result = profiles.merge_ingress_values(case['pkg_values'], case['ingress_values'])

        if not case['expected'](result):
            print(struct.decode(result))
            assert.fail(case['failure_message'])
        end
    end
end

test_merge_ingress_values()
