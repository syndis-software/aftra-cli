import argparse
import json
from collections import defaultdict


def get_subset(data, paths, exclude_paths):
    output = defaultdict(dict)
    for path in paths:
        focus_in = data
        focus_out = output
        steps = path.split(".")[:-1]
        final_key = path.split(".")[-1]

        for step in steps:
            focus_in = focus_in[step]
            out = focus_out.get(step, {})
            focus_out[step] = out
            focus_out = out
        focus_out[final_key] = focus_in[final_key]

    for path in exclude_paths:
        steps = path.split(".")[:-1]
        final_key = path.split(".")[-1]
        focus_out = output
        for step in steps:
            focus_out = focus_out[step]
        del focus_out[final_key]

    return output


def fix_anyof_null(data):
    match data:
        case dict():
            for key, value in data.items():
                if isinstance(value, dict) and "anyOf" in value:
                    any_of_conf = value["anyOf"]
                    if {"type": "null"} in any_of_conf:
                        final = any_of_conf[::]
                        final.remove({"type": "null"})
                        if len(final) == 1:
                            data[key] = final[0]
                        else:
                            data[key]["anyOf"] = final
            for item in data.values():
                fix_anyof_null(item)
        case list():
            for item in data:
                fix_anyof_null(item)


PATHS = [
    "openapi",
    "info",
    "paths./api/companies/{company_pk}/opportunities/.post",
    "paths./api/companies/{parent_pk}/syndis-scans.get",
    "paths./api/companies/{parent_pk}/blobs/upload",
    "paths./api/token/.get",
    "paths./api/integrations/syndis-scan/{scan_name}/config.get",
    "paths./api/integrations/syndis-scan/{scan_name}/logs.post",
    "paths./api/integrations/syndis-scan/{scan_name}/scan.post",
    "components.schemas.MaskedToken",
    "components.schemas.SubmitLogEvent",
    "components.schemas.CreateOpportunity",
    "components.schemas.HTTPValidationError",
    "components.schemas.OpportunityScore",
    "components.schemas.PaginatedEntityCollection_SyndisScanEntity_",
    "components.schemas.SyndisScanTypes",
    "components.schemas.SyndisScanConfig",
    "components.schemas.SyndisScanEntity",
    "components.schemas.SyndisInternalScanEvent_SyndisCISResult_",
    "components.schemas.SyndisInternalScanEvent_SyndisRiskScore_",
    "components.schemas.SyndisRiskScore",
    "components.schemas.SyndisCISResult",
    "components.schemas.ValidationError",
    "components.schemas.Body_Submit_scan_results",
    "components.schemas.BlobUploadInfo",
    "components.schemas.BlobSignedUploadURLResponse",
]

EXCLUDE_PATHS = [
    "components.schemas.SyndisScanEntity.properties.created",
    "components.schemas.SyndisScanEntity.properties.updated",
    "components.schemas.SyndisScanEntity.properties.pk",
    "components.schemas.SyndisScanEntity.properties.sk",
    "components.schemas.SyndisScanEntity.properties.entityType",
]

# test_cases()
if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="SubsetMaker",
        description="Get a subset of a json file based on paths",
    )
    parser.add_argument("filename")  # positional argument

    args = parser.parse_args()

    data = json.loads(open(args.filename).read())
    subset = get_subset(data, PATHS, EXCLUDE_PATHS)
    fix_anyof_null(subset)
    print(json.dumps(subset, indent=4))
