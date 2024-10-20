<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class StaticSiteContent extends Model
{
    use HasFactory;

    protected $hidden = [
        'id',
        'property',
        'created_at',
        'updated_at'
    ];

    protected $fillable = [
        'property',
        'lang_code',
        'content'
    ];
}
