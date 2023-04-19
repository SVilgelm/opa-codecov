package authz
import future.keywords

replace_rule if {
    replace(input.label)
}

replace(label) if {
    label == "test_label"
}
