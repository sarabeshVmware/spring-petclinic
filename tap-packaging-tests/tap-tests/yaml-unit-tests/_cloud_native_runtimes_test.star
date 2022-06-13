load("@ytt:struct", "struct")
load("@ytt:assert", "assert")
load("@ytt:yaml", "yaml")
load("_profiles.star", "profiles")
load("function-to-test.lib.yaml", "collect_values")

def test_shared_ingress_values():
    cases = [
        {
            'short_description': 'propagate shared ingress values',
            'tap_values': struct.encode({
                'profile': 'light',
                'shared': {
                    'ingress_domain': 'example.com',
                },
            }),
            'expected': lambda result: (result.domain_name == 'example.com'),
            'failure_message': 'expected shared ingress values to be passed through to the individual package values',
        },
        {
            'short_description': 'preserve other package values',
            'tap_values': struct.encode({
                'profile': 'light',
                'shared': {
                    'ingress_domain': 'example.com',
                },
                'cnrs': {
                    'someKey': 'someValue'
                },
            }),
            'expected': lambda result: (result.someKey == 'someValue' and
                                        result.domain_name == 'example.com'),
            'failure_message': 'expected individual package values to be preserved',
        },
        {
            'short_description': 'package values take precedence over shared ingress values',
            'tap_values': struct.encode({
                'profile': 'light',
                'shared': {
                    'ingress_domain': 'example.com',
                },
                'cnrs': {
                    'domain_name': 'domain-xyz.com',
                },
            }),
            'expected': lambda result: result.domain_name == 'domain-xyz.com',
            'failure_message': 'expected individual package values to take precedence over shared values',
        },
        {
            'short_description': 'empty shared values',
            'tap_values': struct.encode({
                'profile': 'light',
                'shared': {},
                'cnrs': {
                    'someKey': 'someValue',
                },
            }),
            'expected': lambda result: (result.someKey == 'someValue'),
            'failure_message': 'empty shared values should not change individual package values',
        },
    ]

    for case in cases:
        if 'short_description' in case:
            print('Running:', case['short_description'])
        end
        result = collect_values(case['tap_values'])

        if not case['expected'](result):
            print(struct.decode(result))
            assert.fail(case['failure_message'])
        end
    end
end

test_shared_ingress_values()
