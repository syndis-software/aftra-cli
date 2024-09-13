import pytest
from subset_maker import fix_anyof_null, get_subset


def test_get_subset():
    basic_data = {"a": 1, "b": {"c": 3, "d": 4, "e": {"e1": "1"}}}
    openapi_like_data = {
        "paths": {"/some/url": {"get": {"get-stuff": 1}, "post": {"post-stuff": 2}}}
    }
    for input, paths, expected in [
        (
            basic_data,
            ["a", "b.d"],
            {"a": 1, "b": {"d": 4}},
        ),
        (basic_data, ["b.e"], {"b": {"e": {"e1": "1"}}}),
        (basic_data, ["b.c", "b.d"], {"b": {"c": 3, "d": 4}}),
        (
            openapi_like_data,
            ["paths./some/url.post"],
            {"paths": {"/some/url": {"post": {"post-stuff": 2}}}},
        ),
    ]:
        result = get_subset(input, paths, exclude_paths=[])
        assert result == expected, f"Failure. Actual {result}; Expected {expected}"


@pytest.mark.parametrize(
    ["schema", "expected_schema"],
    [
        (
            {
                "anyOf": [
                    {
                        "type": "string",
                    },
                    {"type": "null"},
                ],
            },
            {"type": "string"},
        ),
        (
            {"type": "string"},
            {"type": "string"},
        ),
    ],
)
def test_fix_anyof_null(schema, expected_schema):
    data = {
        "paths": {
            "/path": {
                "method": {
                    "parameters": [
                        {"schema": schema},
                    ]
                }
            }
        }
    }
    expected = {
        "paths": {
            "/path": {
                "method": {
                    "parameters": [
                        {"schema": expected_schema},
                    ]
                }
            }
        }
    }
    fix_anyof_null(data)
    assert data == expected
